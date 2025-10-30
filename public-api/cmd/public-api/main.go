package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/gops/agent"
	"github.com/rs/zerolog/log"
	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/security"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"github.com/runtime-radar/runtime-radar/lib/server/healthcheck"
	enf_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/auth"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/build"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/client"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/config"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/database"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/server"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/service"
	"go.uber.org/automaxprocs/maxprocs"
	"gorm.io/gorm"
)

const (
	// TLS cert file name.
	certFile = "cert.pem"
	// TLS key file name.
	keyFile = "key.pem"
	// CA cert file name.
	caFile = "ca.pem"

	// timeout on graceful shutdown.
	gracefulTimeout = 15 * time.Second
)

var (
	// Channel for stopping the program.
	shutdown = make(chan struct{})
)

func main() {
	cfg := config.New()
	initLogger(cfg.LogLevel)

	log.Info().Str("build_release", build.Release).Str("build_branch", build.Branch).Str("build_commit", build.Commit).Str("build_date", build.Date).Msgf("-> %s started", build.AppName)
	defer log.Info().Msgf("<- %s exited", build.AppName)

	if _, err := maxprocs.Set(maxprocs.Logger(log.Debug().Msgf)); err != nil {
		log.Warn().Msgf("Can't set maxprocs: %v", err)
	}

	if err := agent.Listen(agent.Options{
		Addr: cfg.GopsAddr,
	}); err != nil {
		log.Fatal().Msgf("### Failed to start gops agent: %v", err)
	}
	defer agent.Close()

	go signalListener()

	var verifier jwt.Verifier
	var err error
	if cfg.Auth {
		verifier, _, err = jwt.NewKeyVerifier(cfg.TokenKey)
		if err != nil {
			log.Fatal().Msgf("### Failed to instantiate key verifier: %v", err)
		}
	}

	// Connect to DB
	db, closeDB, err := database.New(cfg.PostgresAddr, cfg.PostgresDB, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresSSLMode, cfg.PostgresSSLCheckCert)
	if err != nil {
		log.Fatal().Msgf("### Failed to open DB: %v", err)
	}
	defer closeDB()

	// Recreate DB from scratch, or migrate automatically when needed
	if err := database.Migrate(db, cfg.NewDB); err != nil {
		log.Fatal().Msgf("### Failed to migrate DB: %v", err)
	}

	var tlsConfig *tls.Config
	if cfg.TLS {
		// Load TLS config
		tlsConfig, err = security.LoadTLS(caFile, certFile, keyFile)
		if err != nil {
			log.Fatal().Msgf("### Failed to load TLS config: %v", err)
		}
	}

	token, err := hex.DecodeString(cfg.TokenKey)
	if err != nil {
		log.Fatal().Msgf("### Failed to convert tokenKey to []byte: %v", err)
	}

	authAPI, closeAuthAPI, err := client.NewAuthAPI(cfg.AuthAPIURL, tlsConfig, token)
	if err != nil {
		log.Fatal().Msgf("### Failed to connect to AuthAPI: %v", err)
	}
	defer closeAuthAPI()

	ruleController, closeRC, err := client.NewRuleController(cfg.PolicyEnforcerGRPCAddr, tlsConfig, token)
	if err != nil {
		log.Fatal().Msgf("### Failed to connect to Policy Enforcer: %v", err)
	}
	defer closeRC()

	runtimeHistory, closeRH, err := client.NewRuntimeHistory(cfg.HistoryAPIGRPCAddr, tlsConfig, token)
	if err != nil {
		log.Fatal().Msgf("### Failed to connect to History API: %v", err)
	}
	defer closeRH()

	salt, err := hex.DecodeString(cfg.AccessTokenSalt)
	if err != nil {
		log.Fatal().Msgf("### Failed to convert accessTokenSalt to []byte: %v", err)
	} else if len(salt) != service.SaltSizeBytes {
		log.Fatal().Msgf("### Access token salt size must be 64 bytes, got %d", len(salt))
	}

	ruleSvc, accessTokenSvc, runtimeHistorySvc := composeServices(db, authAPI, ruleController, runtimeHistory, salt, cfg.Auth, verifier, token)
	srv := server.New(cfg.ListenHTTPAddr, tlsConfig, accessTokenSvc, ruleSvc, runtimeHistorySvc)

	// Create and Run the instrumentation HTTP server for probes, etc.
	iSrv := server.NewInstrumentation(cfg.InstrumentationAddr)
	go func() {
		if err := iSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("### Can't serve instrumentation HTTP requests: %v", err)
		}
	}()
	log.Info().Msgf("Instrumentation HTTP server listening at %v", cfg.InstrumentationAddr)

	go func() {
		if cfg.TLS {
			if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msgf("### Failed to serve HTTP requests: %v", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msgf("### Failed to serve HTTP requests: %v", err)
			}
		}
	}()

	log.Info().Msgf("HTTP server listening at %v", srv.Addr)

	healthcheck.SetReady() // <-- turn on ready status for k8s

	<-shutdown

	log.Info().Msg("service stopping gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	log.Info().Msg("HTTP server stopping gracefully")
	srv.Shutdown(ctx) // we don't care about errors here

	log.Info().Msg("Instrumentation HTTP server stopping gracefully")
	_ = iSrv.Shutdown(ctx)
}

func signalListener() {
	defer close(shutdown)

	sigTerm := make(chan os.Signal, 10)
	sigIgnore := make(chan os.Signal, 10)

	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigIgnore, syscall.SIGHUP)

	// Wait for signals
	for {
		select {
		case s := <-sigTerm:
			log.Info().Str("signal", s.String()).Msg("Signal caught, terminating")
			return
		case s := <-sigIgnore:
			// Ignoring, like with "nohup"
			log.Info().Str("signal", s.String()).Msg("Signal caught, ignoring")
		}
	}
}
func composeServices(
	db *gorm.DB,
	usersGetter client.UsersGetter,
	ruleController enf_api.RuleControllerClient,
	runtimeHistory history_api.RuntimeHistoryClient,
	accessTokenSalt []byte,
	isAuth bool,
	jwtVerifier jwt.Verifier,
	tokenKey []byte,
) (ruleSvc service.Rule, accessTokenSvc service.AccessToken, runtimeHistorySvc service.RuntimeHistory) {
	accessTokenSvc = &service.AccessTokenGeneric{
		accessTokenSalt,
		tokenKey,
		&database.AccessTokenDatabase{DB: db},
	}

	if isAuth {
		accessTokenSvc = &service.AccessTokenAuth{
			accessTokenSvc,
			tokenKey,
			jwtVerifier,
		}
	}

	ruleSvc = &service.RuleGeneric{ruleController}
	ruleSvc = &service.RuleAuth{
		ruleSvc,
		&auth.Verifier{
			usersGetter,
			&database.AccessTokenDatabase{db},
			accessTokenSalt,
		},
	}

	runtimeHistorySvc = &service.RuntimeHistoryGeneric{runtimeHistory}
	runtimeHistorySvc = &service.RuntimeHistoryAuth{
		runtimeHistorySvc,
		&auth.Verifier{
			usersGetter,
			&database.AccessTokenDatabase{db},
			accessTokenSalt,
		},
	}

	accessTokenSvc = &service.AccessTokenLogging{accessTokenSvc}
	ruleSvc = &service.RuleLogging{ruleSvc}
	runtimeHistorySvc = &service.RuntimeHistoryLogging{runtimeHistorySvc}

	return
}

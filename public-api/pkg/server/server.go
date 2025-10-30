package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/runtime-radar/runtime-radar/lib/server/healthcheck"
	"github.com/runtime-radar/runtime-radar/lib/server/middleware"
	local_middleware "github.com/runtime-radar/runtime-radar/public-api/pkg/server/middleware"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/service"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/service/constructor"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
)

// New constructs and configures new *http.Server capable of serving application endpoints.
func New(httpAddr string, tlsConfig *tls.Config, accessTokenSvc service.AccessToken, ruleSvc service.Rule, runtimeHistorySvc service.RuntimeHistory) *http.Server {
	r := mux.NewRouter()

	return &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         httpAddr,
		Handler:      setupRouter(r, accessTokenSvc, ruleSvc, runtimeHistorySvc),
		TLSConfig:    tlsConfig,
	}
}

func NewInstrumentation(listenAddress string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/ready", healthcheck.ReadyHandler)
	mux.HandleFunc("/live", healthcheck.LiveHandler)

	h := alice.New(
		middleware.Log,
		middleware.Recovery,
	).Then(mux)

	return &http.Server{
		ReadTimeout: readTimeout,
		Addr:        listenAddress,
		Handler:     h,
	}
}

func setupRouter(r *mux.Router, accessTokenSvc service.AccessToken, ruleSvc service.Rule, runtimeHistorySvc service.RuntimeHistory) http.Handler {
	r.StrictSlash(true)

	corsOpts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization", "X-Auth-Key"},
	}

	h := alice.New(
		middleware.Log,
		middleware.Recovery,
		cors.New(corsOpts).Handler,
		local_middleware.JWT,
		local_middleware.AccessToken,
		local_middleware.Correlation,
	).Then(r)

	r.Handle("/api/v1/access-token", constructor.AccessTokenCreate(accessTokenSvc)).Methods(http.MethodPost)
	r.Handle("/api/v1/access-token/page/{page_num:[0-9]+}", constructor.AccessTokenListPage(accessTokenSvc)).Methods(http.MethodGet)
	r.Handle("/api/v1/access-token/{id}", constructor.AccessTokenDelete(accessTokenSvc)).Methods(http.MethodDelete)
	r.Handle("/api/v1/access-token/{id}", constructor.AccessTokenGetByID(accessTokenSvc)).Methods(http.MethodGet)
	r.Handle("/api/v1/access-token/invalidate-access-tokens", constructor.AccessTokenInvalidateAll(accessTokenSvc)).Methods(http.MethodPost)

	// policy-enforcer
	r.Handle("/api/v1/public-api/rule", constructor.RuleCreate(ruleSvc)).Methods(http.MethodPost)
	r.Handle("/api/v1/public-api/rule/page/{page_num:[0-9]+}", constructor.RuleListPage(ruleSvc)).Methods(http.MethodGet)
	r.Handle("/api/v1/public-api/rule/notify-targets-in-use", constructor.RuleNotifyTargetsInUse(ruleSvc)).Methods(http.MethodGet)
	r.Handle("/api/v1/public-api/rule/{id}", constructor.RuleRead(ruleSvc)).Methods(http.MethodGet)
	r.Handle("/api/v1/public-api/rule/{id}", constructor.RuleUpdate(ruleSvc)).Methods(http.MethodPatch)
	r.Handle("/api/v1/public-api/rule/{id}", constructor.RuleDelete(ruleSvc)).Methods(http.MethodDelete)

	// history-api
	r.Handle("/api/v1/public-api/runtime-event/slice/{direction:left|right}", constructor.RuntimeHistoryListEventsSlice(runtimeHistorySvc)).
		Methods(http.MethodGet).
		Queries("cursor", `{cursor:[a-zA-Z0-9\-:.]+}`)

	return h
}

package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/runtime-radar/runtime-radar/lib/server/healthcheck"
	"github.com/runtime-radar/runtime-radar/lib/server/middleware"
	"github.com/runtime-radar/runtime-radar/runtime-monitor/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	readTimeout  = 2 * time.Minute
	writeTimeout = 2 * time.Minute
	// Maximum message size for grpc request
	MaxRecvMsgSize = 10 * 1024 * 1024 // 10MB
)

// New constructs and configures new *http.Server capable of serving application and gRPC gateway endpoints.
func New(httpAddr, grpcAddr string, tlsConfig *tls.Config) (*http.Server, error) {
	mux := http.NewServeMux()
	gwMux, err := newGWMux(context.Background(), grpcAddr, tlsConfig)
	if err != nil {
		return nil, err
	}

	h := setupRouter(mux, gwMux)

	s := &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         httpAddr,
		Handler:      h,
		TLSConfig:    tlsConfig,
	}

	return s, nil
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
		Addr:    listenAddress,
		Handler: h,
	}
}

func setupRouter(mux *http.ServeMux, gwMux *runtime.ServeMux) http.Handler {
	mux.Handle("/", gwMux)

	corsOpts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
	}

	h := alice.New(
		middleware.Log,
		middleware.Recovery,
		cors.New(corsOpts).Handler,
	).Then(mux)

	return h
}

func newGWMux(ctx context.Context, grpcAddr string, tlsConfig *tls.Config) (*runtime.ServeMux, error) {
	m := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		// Configure to always output same snake_case field names as in .*proto
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	var creds credentials.TransportCredentials
	if tlsConfig != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize))}
	if err := api.RegisterConfigControllerHandlerFromEndpoint(ctx, m, grpcAddr, opts); err != nil {
		return nil, err
	}

	return m, nil
}

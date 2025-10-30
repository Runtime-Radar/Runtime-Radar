package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"google.golang.org/grpc/metadata"
)

type accessTokenContextKeyType string

const (
	accessTokenContextKey accessTokenContextKeyType = "access_token"
)

const (
	AuthorizationHeader = "Authorization"
	AccessTokenHeader   = "Access-Token"
	CorrelationHeader   = "Grpc-Metadata-Correlation-Id"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if jwtHeader := r.Header.Get(AuthorizationHeader); jwtHeader != "" {
			ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", jwtHeader))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if accessTokenHeader := r.Header.Get(AccessTokenHeader); accessTokenHeader != "" {
			ctx = context.WithValue(ctx, accessTokenContextKey, accessTokenHeader)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get(CorrelationHeader)
		if corrID == "" {
			corrID = uuid.New().String()
		}

		ctx := r.Context()
		ctx = metadata.AppendToOutgoingContext(ctx, interceptor.CorrelationHeader, corrID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(accessTokenContextKey).(string)
	return token, ok
}

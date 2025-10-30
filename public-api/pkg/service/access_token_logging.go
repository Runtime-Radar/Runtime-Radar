package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/model"
)

type AccessTokenLogging struct {
	AccessToken
}

func (a *AccessTokenLogging) Create(ctx context.Context, req *model.CreateAccessTokenReq) (id uuid.UUID, token string, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Str("access_token", "HIDDEN").
			Stringer("token_id", id).
			Msg("Called AccessToken.Create")
	}(time.Now())

	id, token, err = a.AccessToken.Create(ctx, req)
	return
}

func (a *AccessTokenLogging) ListPage(ctx context.Context, pageNum, pageSize int, order string) (ts []*model.AccessTokenResp, total int, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Int("page_num", pageNum).
			Int("page_size", pageSize).
			Int("total", total).
			Int("tokens", len(ts)).
			Msg("Called AccessToken.ListPage")
	}(time.Now())

	ts, total, err = a.AccessToken.ListPage(ctx, pageNum, pageSize, order)
	return
}

func (a *AccessTokenLogging) GetByID(ctx context.Context, id uuid.UUID) (at *model.AccessTokenResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Stringer("id", id).
			Interface("at", at).
			Msg("Called AccessToken.GetByID")
	}(time.Now())

	at, err = a.AccessToken.GetByID(ctx, id)
	return
}

func (a *AccessTokenLogging) Delete(ctx context.Context, id uuid.UUID) (err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Stringer("id", id).
			Msg("Called AccessToken.Delete")
	}(time.Now())

	err = a.AccessToken.Delete(ctx, id)
	return
}

func (a *AccessTokenLogging) InvalidateAll(ctx context.Context) (err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Msg("Called AccessToken.InvalidateAll")
	}(time.Now())

	err = a.AccessToken.InvalidateAll(ctx)
	return
}

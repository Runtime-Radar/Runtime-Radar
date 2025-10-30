package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RuleLogging struct {
	api.RuleControllerServer
}

func (rl *RuleLogging) Create(ctx context.Context, req *api.Rule) (resp *api.CreateRuleResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Bool("audit", true).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.Create")
	}(time.Now())

	resp, err = rl.RuleControllerServer.Create(ctx, req)
	return
}

func (rl *RuleLogging) Read(ctx context.Context, req *api.ReadRuleReq) (resp *api.ReadRuleResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.Read")
	}(time.Now())

	resp, err = rl.RuleControllerServer.Read(ctx, req)
	return
}

func (rl *RuleLogging) Update(ctx context.Context, req *api.Rule) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Bool("audit", true).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.Update")
	}(time.Now())

	resp, err = rl.RuleControllerServer.Update(ctx, req)
	return
}

func (rl *RuleLogging) Delete(ctx context.Context, req *api.DeleteRuleReq) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Bool("audit", true).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.Delete")
	}(time.Now())

	resp, err = rl.RuleControllerServer.Delete(ctx, req)
	return
}

func (rl *RuleLogging) ListPage(ctx context.Context, req *api.ListRulePageReq) (resp *api.ListRulePageResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetRules())).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.ListPage")
	}(time.Now())

	resp, err = rl.RuleControllerServer.ListPage(ctx, req)
	return
}

func (rl *RuleLogging) NotifyTargetsInUse(ctx context.Context, req *api.NotifyTargetsInUseReq) (resp *api.NotifyTargetsInUseResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuleControllerServer.NotifyTargetsInUse")
	}(time.Now())

	resp, err = rl.RuleControllerServer.NotifyTargetsInUse(ctx, req)
	return
}

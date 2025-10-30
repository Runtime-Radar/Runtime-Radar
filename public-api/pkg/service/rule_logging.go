package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	enf_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RuleLogging struct {
	enf_api.RuleControllerClient
}

func (r *RuleLogging) Create(ctx context.Context, req *enf_api.Rule, opts ...grpc.CallOption) (resp *enf_api.CreateRuleResp, err error) {
	defer func(t0 time.Time) {
		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Msg("Called Rule.Create")
	}(time.Now())

	resp, err = r.RuleControllerClient.Create(ctx, req, opts...)
	return
}

func (r *RuleLogging) Read(ctx context.Context, req *enf_api.ReadRuleReq, opts ...grpc.CallOption) (resp *enf_api.ReadRuleResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Interface("result", resp).
			Msg("Called Rule.Read")
	}(time.Now())

	resp, err = r.RuleControllerClient.Read(ctx, req, opts...)
	return
}

func (r *RuleLogging) Update(ctx context.Context, req *enf_api.Rule, opts ...grpc.CallOption) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Interface("result", resp).
			Msg("Called Rule.Update")
	}(time.Now())

	resp, err = r.RuleControllerClient.Update(ctx, req, opts...)
	return
}

func (r *RuleLogging) Delete(ctx context.Context, req *enf_api.DeleteRuleReq, opts ...grpc.CallOption) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Interface("result", resp).
			Msg("Called Rule.Delete")
	}(time.Now())

	resp, err = r.RuleControllerClient.Delete(ctx, req, opts...)
	return
}

func (r *RuleLogging) ListPage(ctx context.Context, req *enf_api.ListRulePageReq, opts ...grpc.CallOption) (resp *enf_api.ListRulePageResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Interface("result", len(resp.GetRules())).
			Msg("Called Rule.ListPage")
	}(time.Now())

	resp, err = r.RuleControllerClient.ListPage(ctx, req, opts...)
	return
}

func (r *RuleLogging) NotifyTargetsInUse(ctx context.Context, req *enf_api.NotifyTargetsInUseReq, opts ...grpc.CallOption) (resp *enf_api.NotifyTargetsInUseResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Stringer("correlation_id", corrID).
			Interface("args", req).
			Interface("result", resp).
			Msg("Called Rule.NotifyTargetsInUse")
	}(time.Now())

	resp, err = r.RuleControllerClient.NotifyTargetsInUse(ctx, req, opts...)
	return
}

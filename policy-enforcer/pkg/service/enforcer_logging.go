package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/api"
)

type EnforcerLogging struct {
	api.EnforcerServer
}

func (el *EnforcerLogging) EvaluatePolicyRuntimeEvent(ctx context.Context, req *api.EvaluatePolicyRuntimeEventReq) (resp *api.EvaluatePolicyRuntimeEventReq, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called EnforcerServer.EvaluatePolicyRuntimeEvent")
	}(time.Now())

	resp, err = el.EnforcerServer.EvaluatePolicyRuntimeEvent(ctx, req)
	return
}

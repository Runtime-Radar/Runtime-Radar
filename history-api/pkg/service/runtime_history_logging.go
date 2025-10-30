package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	processor_api "github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
)

type RuntimeHistoryLogging struct {
	api.RuntimeHistoryServer
}

func (sl *RuntimeHistoryLogging) Read(ctx context.Context, req *api.ReadRuntimeEventReq) (resp *processor_api.RuntimeEvent, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RuntimeHistoryServer.Read")
	}(time.Now())

	resp, err = sl.RuntimeHistoryServer.Read(ctx, req)
	return
}

func (sl *RuntimeHistoryLogging) ListRuntimeEventSlice(ctx context.Context, req *api.ListRuntimeEventSliceReq) (resp *api.ListRuntimeEventSliceResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetRuntimeEvents())).
			Stringer("correlation_id", corrID).
			Msg("Called RuntimeHistoryServer.ListRuntimeEventSlice")
	}(time.Now())

	resp, err = sl.RuntimeHistoryServer.ListRuntimeEventSlice(ctx, req)
	return
}

func (sl *RuntimeHistoryLogging) FilterRuntimeEventSlice(ctx context.Context, req *api.FilterRuntimeEventSliceReq) (resp *api.ListRuntimeEventSliceResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetRuntimeEvents())).
			Stringer("correlation_id", corrID).
			Msg("Called RuntimeHistoryServer.FilterRuntimeEventSlice")
	}(time.Now())

	resp, err = sl.RuntimeHistoryServer.FilterRuntimeEventSlice(ctx, req)
	return
}

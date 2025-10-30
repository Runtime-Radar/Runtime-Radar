package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"google.golang.org/grpc"
)

type RuntimeHistoryLogging struct {
	history_api.RuntimeHistoryClient
}

func (rhl *RuntimeHistoryLogging) ListRuntimeEventSlice(ctx context.Context, req *history_api.ListRuntimeEventSliceReq, opts ...grpc.CallOption) (resp *history_api.ListRuntimeEventSliceResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).
			Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetRuntimeEvents())).
			Stringer("correlation_id", corrID).
			Msg("Called RuntimeHistory.ListRuntimeEventSlice")
	}(time.Now())

	resp, err = rhl.RuntimeHistoryClient.ListRuntimeEventSlice(ctx, req, opts...)
	return
}

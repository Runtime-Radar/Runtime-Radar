package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DetectorLogging struct {
	api.DetectorControllerServer
}

func (dl *DetectorLogging) Create(ctx context.Context, req *api.CreateDetectorReq) (resp *api.CreateDetectorResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)
		detectorSize := base64.StdEncoding.DecodedLen(len(req.GetWasmBase64())) / 1024

		log.Err(err).Str("delay", time.Since(t0).String()).
			Str("args", fmt.Sprintf("detector size %d kB", detectorSize)).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called DetectorControllerServer.Create")
	}(time.Now())

	resp, err = dl.DetectorControllerServer.Create(ctx, req)
	return
}

func (dl *DetectorLogging) Delete(ctx context.Context, req *api.DeleteDetectorReq) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called DetectorControllerServer.Delete")
	}(time.Now())

	resp, err = dl.DetectorControllerServer.Delete(ctx, req)
	return
}

func (dl *DetectorLogging) ListPage(ctx context.Context, req *api.ListDetectorPageReq) (resp *api.ListDetectorPageResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetDetectors())).
			Stringer("correlation_id", corrID).
			Msg("Called DetectorControllerServer.ListPage")
	}(time.Now())

	resp, err = dl.DetectorControllerServer.ListPage(ctx, req)
	return
}

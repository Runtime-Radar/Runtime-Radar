package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/auth-center/api"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleLogging struct {
	api.RoleControllerServer
}

func (rl *RoleLogging) Read(ctx context.Context, req *api.ReadRoleReq) (resp *api.ReadRoleResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RoleServer.Read")
	}(time.Now())

	resp, err = rl.RoleControllerServer.Read(ctx, req)
	return
}

func (rl *RoleLogging) ReadList(ctx context.Context, _ *emptypb.Empty) (resp *api.ReadListRoleResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called RoleServer.ReadList")
	}(time.Now())

	resp, err = rl.RoleControllerServer.ReadList(ctx, &emptypb.Empty{})
	return
}

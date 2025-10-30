package service

import (
	"context"

	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"google.golang.org/grpc"
)

type RuntimeHistoryAuth struct {
	history_api.RuntimeHistoryClient

	Verifier jwt.Verifier
}

func (rha *RuntimeHistoryAuth) ListRuntimeEventSlice(ctx context.Context, req *history_api.ListRuntimeEventSliceReq, opts ...grpc.CallOption) (*history_api.ListRuntimeEventSliceResp, error) {
	if err := rha.Verifier.VerifyPermission(ctx, jwt.PermissionEvents, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}

	return rha.RuntimeHistoryClient.ListRuntimeEventSlice(ctx, req, opts...)
}

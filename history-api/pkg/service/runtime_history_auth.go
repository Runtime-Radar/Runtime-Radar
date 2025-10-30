package service

import (
	"context"

	processor_api "github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

// RuntimeHistoryAuth is a layer for jwt-based authentication.
// Base server interface should not be embedded here unlike
// in implementations of other layers.
// All required methods should be explicitly implemented to ensure
// that new methods of the basic server are implemented for auth layer.
type RuntimeHistoryAuth struct {
	// UnsafeRuntimeHistoryServer is embedded to opt out of forward
	// compatibility promised by protobuf library.
	// It merely contains an empty `mustEmbedUnimplementedStatsServer()`
	// method.
	api.UnsafeRuntimeHistoryServer

	// RuntimeHistoryServer is a base server interface to pass
	// response to the next layer.
	RuntimeHistoryServer api.RuntimeHistoryServer
	Verifier             jwt.Verifier
}

func (rha *RuntimeHistoryAuth) Read(ctx context.Context, req *api.ReadRuntimeEventReq) (resp *processor_api.RuntimeEvent, err error) {
	if err := rha.Verifier.VerifyPermission(ctx, jwt.PermissionEvents, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	resp, err = rha.RuntimeHistoryServer.Read(ctx, req)
	return
}

func (rha *RuntimeHistoryAuth) ListRuntimeEventSlice(ctx context.Context, req *api.ListRuntimeEventSliceReq) (resp *api.ListRuntimeEventSliceResp, err error) {
	if err := rha.Verifier.VerifyPermission(ctx, jwt.PermissionEvents, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	resp, err = rha.RuntimeHistoryServer.ListRuntimeEventSlice(ctx, req)
	return
}

func (rha *RuntimeHistoryAuth) FilterRuntimeEventSlice(ctx context.Context, req *api.FilterRuntimeEventSliceReq) (resp *api.ListRuntimeEventSliceResp, err error) {
	if err := rha.Verifier.VerifyPermission(ctx, jwt.PermissionEvents, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	resp, err = rha.RuntimeHistoryServer.FilterRuntimeEventSlice(ctx, req)
	return
}

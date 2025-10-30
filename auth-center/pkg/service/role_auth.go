package service

import (
	"context"

	"github.com/runtime-radar/runtime-radar/auth-center/api"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleAuth struct {
	api.UnsafeRoleControllerServer
	RoleControllerServer api.RoleControllerServer
	Verifier             jwt.Verifier
}

func (ra *RoleAuth) Read(ctx context.Context, req *api.ReadRoleReq) (resp *api.ReadRoleResp, err error) {
	if err := ra.Verifier.VerifyPermission(ctx, jwt.PermissionRoles, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	resp, err = ra.RoleControllerServer.Read(ctx, req)
	return
}

func (ra *RoleAuth) ReadList(ctx context.Context, _ *emptypb.Empty) (resp *api.ReadListRoleResp, err error) {
	if err := ra.Verifier.VerifyPermission(ctx, jwt.PermissionRoles, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	resp, err = ra.RoleControllerServer.ReadList(ctx, &emptypb.Empty{})
	return
}

package service

import (
	"context"

	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"github.com/runtime-radar/runtime-radar/notifier/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IntegrationAuth struct {
	api.UnsafeIntegrationControllerServer

	IntegrationControllerServer api.IntegrationControllerServer
	Verifier                    jwt.Verifier
}

func (ia *IntegrationAuth) Create(ctx context.Context, req *api.Integration) (*api.CreateIntegrationResp, error) {
	if err := ia.Verifier.VerifyPermission(ctx, jwt.PermissionIntegrations, jwt.ActionCreate); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	return ia.IntegrationControllerServer.Create(ctx, req)
}

func (ia *IntegrationAuth) Read(ctx context.Context, req *api.ReadIntegrationReq) (*api.Integration, error) {
	if err := ia.Verifier.VerifyPermission(ctx, jwt.PermissionIntegrations, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	return ia.IntegrationControllerServer.Read(ctx, req)
}

func (ia *IntegrationAuth) Update(ctx context.Context, req *api.Integration) (*emptypb.Empty, error) {
	if err := ia.Verifier.VerifyPermission(ctx, jwt.PermissionIntegrations, jwt.ActionUpdate); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	return ia.IntegrationControllerServer.Update(ctx, req)
}

func (ia *IntegrationAuth) Delete(ctx context.Context, req *api.DeleteIntegrationReq) (*emptypb.Empty, error) {
	if err := ia.Verifier.VerifyPermission(ctx, jwt.PermissionIntegrations, jwt.ActionDelete); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	return ia.IntegrationControllerServer.Delete(ctx, req)
}

func (ia *IntegrationAuth) List(ctx context.Context, req *api.ListIntegrationReq) (*api.ListIntegrationResp, error) {
	if err := ia.Verifier.VerifyPermission(ctx, jwt.PermissionIntegrations, jwt.ActionRead); err != nil {
		return nil, errcommon.PermissionErrorToStatus(err)
	}
	return ia.IntegrationControllerServer.List(ctx, req)
}

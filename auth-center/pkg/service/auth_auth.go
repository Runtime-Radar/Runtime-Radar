package service

import (
	"context"

	"github.com/runtime-radar/runtime-radar/auth-center/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthAuth struct {
	api.UnsafeAuthControllerServer
	AuthControllerServer api.AuthControllerServer
}

func (aa *AuthAuth) SignIn(ctx context.Context, req *api.SignInReq) (resp *api.SignInResp, err error) {
	return aa.AuthControllerServer.SignIn(ctx, req)
}

func (aa *AuthAuth) RefreshTokens(ctx context.Context, empty *emptypb.Empty) (resp *api.SignInResp, err error) {
	return aa.AuthControllerServer.RefreshTokens(ctx, empty)
}

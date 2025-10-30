package service

import (
	"context"
	"net/url"

	"github.com/runtime-radar/runtime-radar/cs-manager/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InfoGeneric struct {
	api.UnimplementedInfoControllerServer
	Version      string
	CentralCSURL *url.URL
}

func (ig *InfoGeneric) GetVersion(_ context.Context, _ *emptypb.Empty) (*api.Version, error) {
	return &api.Version{
		Version: ig.Version,
	}, nil
}

func (ig *InfoGeneric) GetCentralCSURL(_ context.Context, _ *emptypb.Empty) (*api.URL, error) {
	return &api.URL{
		Url: ig.CentralCSURL.String(),
	}, nil
}

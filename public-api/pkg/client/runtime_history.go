package client

import (
	"crypto/tls"

	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/build"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRuntimeHistory(address string, tlsConfig *tls.Config, tokenKey []byte) (history_api.RuntimeHistoryClient, func() error, error) {
	var creds credentials.TransportCredentials
	if tlsConfig != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	if len(tokenKey) > 0 && tlsConfig != nil {
		rp := &jwt.RolePermissions{
			Events: &jwt.Permission{
				Actions: []jwt.Action{jwt.ActionRead},
			},
		}

		creds, err := jwt.GeneratePerRPCCredentials(tokenKey, build.AppName, rp)
		if err != nil {
			return nil, nil, err
		}

		opts = append(opts, grpc.WithPerRPCCredentials(creds))
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, nil, err
	}

	client := history_api.NewRuntimeHistoryClient(conn)

	return client, conn.Close, nil
}

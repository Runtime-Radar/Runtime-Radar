package client

import (
	"crypto/tls"

	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/build"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/server"
	enforcer_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRuleController(address string, tlsConfig *tls.Config, tokenKey []byte) (enforcer_api.RuleControllerClient, func() error, error) {
	var creds credentials.TransportCredentials
	if tlsConfig != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(server.MaxRecvMsgSize)),
		grpc.WithNoProxy(), // This client is used for internal interactions where proxy settings are never required.
	}

	if len(tokenKey) > 0 {
		rp := &jwt.RolePermissions{
			Rules: &jwt.Permission{
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

	client := enforcer_api.NewRuleControllerClient(conn)

	return client, conn.Close, nil
}

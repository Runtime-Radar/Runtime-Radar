package client

import (
	"crypto/tls"

	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	enf_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/build"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRuleController(address string, tlsConfig *tls.Config, tokenKey []byte) (enf_api.RuleControllerClient, func() error, error) {
	var creds credentials.TransportCredentials
	if tlsConfig != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	if len(tokenKey) > 0 && tlsConfig != nil {
		rp := &jwt.RolePermissions{
			Rules: &jwt.Permission{
				Actions: []jwt.Action{
					jwt.ActionCreate, jwt.ActionRead, jwt.ActionUpdate, jwt.ActionDelete,
				},
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

	client := enf_api.NewRuleControllerClient(conn)

	return client, conn.Close, nil
}

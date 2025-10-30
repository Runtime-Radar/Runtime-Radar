package client

import (
	"crypto/tls"

	cluster_api "github.com/runtime-radar/runtime-radar/cluster-manager/api"
	"github.com/runtime-radar/runtime-radar/cs-manager/pkg/build"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClusterController(address string, tlsConfig *tls.Config, tokenKey []byte) (cluster_api.ClusterControllerClient, func() error, error) {
	var creds credentials.TransportCredentials
	if tlsConfig != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	if len(tokenKey) > 0 {
		rp := &jwt.RolePermissions{
			Clusters: &jwt.Permission{
				Actions: []jwt.Action{
					jwt.ActionRead,
					jwt.ActionUpdate,
					jwt.ActionDelete,
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

	client := cluster_api.NewClusterControllerClient(conn)

	return client, conn.Close, nil
}

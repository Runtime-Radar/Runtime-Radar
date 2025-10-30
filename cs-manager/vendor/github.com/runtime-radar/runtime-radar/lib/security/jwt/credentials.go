package jwt

import (
	"context"

	"google.golang.org/grpc/credentials"
)

type creds struct {
	signed string
}

// GetRequestMetadata is a method to satisfy credentials.PerRPCCredentials so token could be used in grpc.WithPerRPCCredentials.
func (c creds) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + c.signed,
	}, nil
}

// RequireTransportSecurity is a method to satisfy credentials.PerRPCCredentials so token could be used in grpc.WithPerRPCCredentials.
func (c creds) RequireTransportSecurity() bool {
	return true
}

func GeneratePerRPCCredentials(key []byte, name string, rolePermissions *RolePermissions) (credentials.PerRPCCredentials, error) {
	token, err := GenerateServiceToken(key, name, rolePermissions)
	if err != nil {
		return nil, err
	}

	return creds{token}, nil
}

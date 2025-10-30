//go:build tools

// Package tools is needed to make project stick to particular versions of tools by putting version data
// in go.mod so that running "go install ..." without tags will generate needed binaries with expected behavior.
package tools

//nolint: typecheck
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/gops"
	_ "golang.org/x/vuln/cmd/govulncheck"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

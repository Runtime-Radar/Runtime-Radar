//go:build tools

// Package tools is needed to make project stick to particular versions of tools by putting version data
// in go.mod so that running "go install ..." without tags will generate needed binaries with expected behavior.
package tools

//nolint: typecheck
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/vuln/cmd/govulncheck"
)

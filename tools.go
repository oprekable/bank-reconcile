//go:build tools
// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/vektra/mockery/v2"
	_ "golang.org/x/tools/cmd/deadcode"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
)

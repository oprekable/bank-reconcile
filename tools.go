//go:build tools
// +build tools

package main

import (
	// Import syntax required to install tool golangci-lint by make script
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// Import syntax required to install tool wire by make script
	_ "github.com/google/wire/cmd/wire"
	// Import syntax required to install tool mockery by make script
	_ "github.com/vektra/mockery/v2"
	// Import syntax required to install tool deadcode by make script
	_ "golang.org/x/tools/cmd/deadcode"
	// Import syntax required to install tool goimports by make script
	_ "golang.org/x/tools/cmd/goimports"
	// Import syntax required to install tool fieldalignment by make script
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	// Import syntax required to install tool govulncheck by make script
	_ "golang.org/x/vuln/cmd/govulncheck"
	// Import syntax required to install tool staticcheck by make script
	_ "honnef.co/go/tools/cmd/staticcheck"
)

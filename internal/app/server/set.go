package server

import (
	"github.com/oprekable/bank-reconcile/internal/app/server/cli"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	cli.Set,
	wire.Bind(new(CliServer), new(*cli.Cli)),
	NewServer,
)

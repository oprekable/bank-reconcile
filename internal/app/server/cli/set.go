package cli

import (
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	hcli.Set,
	NewCli,
)

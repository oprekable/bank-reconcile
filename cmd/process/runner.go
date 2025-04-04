package process

import (
	"github.com/oprekable/bank-reconcile/cmd/helper"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"

	"github.com/spf13/cobra"
)

func Runner(cmd *cobra.Command, args []string) (er error) {
	dBPath := csqlite.DBPath{}

	if root.FlagIsDebugValue {
		dBPath.WriteDBPath = "./reconciliation.db"
	}

	return helper.RunnerSubCommand(cmd, args, dBPath)
}

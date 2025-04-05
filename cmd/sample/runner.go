package sample

import (
	"github.com/oprekable/bank-reconcile/cmd/helper"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/err"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

var wireApp = inject.WireApp

func Runner(cmd *cobra.Command, args []string) (er error) {
	dBPath := csqlite.DBPath{}

	if root.FlagIsDebugValue {
		dBPath.ReadDBPath = "./sample.db"
	}

	r := helper.NewRunner(
		wireApp,
		cmd,
		args,
	)

	return r.Run(
		root.EmbedFS,
		variable.AppName,
		root.FlagTZValue,
		err.RegisteredErrorType,
		root.FlagIsVerboseValue,
		dBPath,
	)
}

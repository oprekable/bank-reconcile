package main

import (
	"embed"
	"fmt"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/cmd/version"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/variable"
	"io"
	"os"
)

//go:embed all:embeds
var embedFS embed.FS

func main() {
	var outPutWriter io.Writer = os.Stdout
	var errWriter io.Writer = os.Stderr

	var subCommands = []cmd.Cmd{
		version.NewCommand(outPutWriter, errWriter),
		sample.NewCommand(variable.AppName, inject.WireAppFn, &embedFS, outPutWriter, errWriter),
		process.NewCommand(variable.AppName, inject.WireAppFn, &embedFS, outPutWriter, errWriter),
	}

	c := root.NewCommand(inject.WireAppFn, &embedFS, outPutWriter, errWriter, subCommands...).
		Init(
			&cmd.MetaData{
				Usage: variable.AppName,
				Short: variable.AppDescShort,
				Long:  variable.AppDescLong,
			},
		)

	if err := c.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

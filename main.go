package main

import (
	"embed"
	"io"
	"os"
	"unsafe"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/cmd/version"
	"github.com/oprekable/bank-reconcile/internal/_inject"
	"github.com/oprekable/bank-reconcile/variable"
)

//go:embed all:embeds
var embedFS embed.FS
var exitFunc = os.Exit

func main() {
	exitFunc(mainLogic(os.Stdout, os.Stderr))
}

func mainLogic(outPutWriter io.Writer, errWriter io.Writer) int {
	var subCommands = []cmd.Cmd{
		version.NewCommand(outPutWriter, errWriter),
		sample.NewCommand(variable.AppName, _inject.WireAppFn, &embedFS, outPutWriter, errWriter),
		process.NewCommand(variable.AppName, _inject.WireAppFn, &embedFS, outPutWriter, errWriter),
	}

	c := root.NewCommand(outPutWriter, errWriter, subCommands...).
		Init(
			&cmd.MetaData{
				Usage: variable.AppName,
				Short: variable.AppDescShort,
				Long:  variable.AppDescLong,
			},
		)

	isHaveErr := c.Execute() != nil
	// Bool to 0 or 1
	return int(*(*byte)(unsafe.Pointer(&isHaveErr)))
}

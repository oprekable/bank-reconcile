package root

import (
	"embed"
	"fmt"
	"io"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/spf13/cobra"
)

type CmdRoot struct {
	c            *cobra.Command
	wireApp      inject.Fn
	embedFS      *embed.FS
	outPutWriter io.Writer
	errWriter    io.Writer
	subCommands  []cmd.Cmd
}

var _ cmd.Cmd = (*CmdRoot)(nil)

func NewCommand(wireApp inject.Fn, embedFS *embed.FS, outPutWriter io.Writer, errWriter io.Writer, subCommands ...cmd.Cmd) *CmdRoot {
	return &CmdRoot{
		c: &cobra.Command{
			Args:          cobra.NoArgs,
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		wireApp:      wireApp,
		embedFS:      embedFS,
		outPutWriter: outPutWriter,
		errWriter:    errWriter,
		subCommands:  subCommands,
	}
}

func (c *CmdRoot) Init(metaData *cmd.MetaData) *cobra.Command {
	c.c.Use = metaData.Usage
	c.c.Short = metaData.Short
	c.c.Long = metaData.Long
	c.c.Example = fmt.Sprintf(
		"%s\n%s\n",
		fmt.Sprintf("Generate sample \n\t%s %s", metaData.Usage, sample.Example),
		fmt.Sprintf("Process data \n\t%s %s", metaData.Usage, process.Example),
	)
	c.c.PersistentPreRunE = c.PersistentPreRunner
	c.c.RunE = c.Runner
	c.c.SetOut(c.outPutWriter)
	c.c.SetErr(c.errWriter)

	for i := range c.subCommands {
		c.c.AddCommand(
			c.subCommands[i].Init(metaData),
		)
	}
	return c.c
}

func (c *CmdRoot) Runner(_ *cobra.Command, _ []string) (er error) {
	return c.c.Help()
}

func (c *CmdRoot) PersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	return nil
}

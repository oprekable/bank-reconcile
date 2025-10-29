package root

import (
	"io"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/spf13/cobra"
)

func errFlagsFunc(command *cobra.Command, err error) error {
	command.Println(err.Error() + "\n")
	command.Println(command.UsageString())
	return err
}

type CmdRoot struct {
	c            *cobra.Command
	outPutWriter io.Writer
	errWriter    io.Writer
	subCommands  []cmd.Cmd
}

var _ cmd.Cmd = (*CmdRoot)(nil)

func NewCommand(outPutWriter io.Writer, errWriter io.Writer, subCommands ...cmd.Cmd) *CmdRoot {
	return &CmdRoot{
		c: &cobra.Command{
			Args:          cobra.NoArgs,
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		outPutWriter: outPutWriter,
		errWriter:    errWriter,
		subCommands:  subCommands,
	}
}

func (c *CmdRoot) Init(metaData *cmd.MetaData) *cobra.Command {
	c.c.Use = metaData.Usage
	c.c.Short = metaData.Short
	c.c.Long = metaData.Long

	c.c.PersistentPreRunE = c.PersistentPreRunner
	c.c.RunE = c.Runner
	c.c.SetOut(c.outPutWriter)
	c.c.SetErr(c.errWriter)

	for i := range c.subCommands {
		c.c.AddCommand(
			c.subCommands[i].Init(metaData),
		)

		c.c.Example += c.subCommands[i].Example()
	}

	c.c.SetFlagErrorFunc(errFlagsFunc)

	return c.c
}

func (c *CmdRoot) Runner(_ *cobra.Command, _ []string) (er error) {
	return c.c.Help()
}

func (c *CmdRoot) PersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	return nil
}

func (c *CmdRoot) Example() string {
	return c.c.Example
}

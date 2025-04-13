package version

import (
	"fmt"
	"io"

	"github.com/oprekable/bank-reconcile/cmd"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/versionhelper"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

type CmdVersion struct {
	c            *cobra.Command
	outPutWriter io.Writer
	errWriter    io.Writer
}

var _ cmd.Cmd = (*CmdVersion)(nil)

func NewCommand(outPutWriter io.Writer, errWriter io.Writer) *CmdVersion {
	return &CmdVersion{
		c: &cobra.Command{
			Args:          cobra.NoArgs,
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		outPutWriter: outPutWriter,
		errWriter:    errWriter,
	}
}

func (c *CmdVersion) Init(_ *cmd.MetaData) *cobra.Command {
	c.c.Use = Usage
	c.c.Aliases = Aliases
	c.c.Short = Short
	c.c.Long = Long
	c.c.RunE = c.Runner
	c.c.SetOut(c.outPutWriter)
	c.c.SetErr(c.errWriter)

	return c.c
}

func (c *CmdVersion) Runner(_ *cobra.Command, _ []string) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	atexit.Add(
		func() {
			_, _ = fmt.Fprintln(c.outPutWriter, "\n-#-")
		},
	)

	version := versionhelper.GetVersion(
		variable.Version,
		variable.BuildDate,
		variable.GitCommit,
		variable.Environment,
	)

	_, _ = fmt.Fprintln(c.outPutWriter, "App\t\t:", variable.AppName)
	_, _ = fmt.Fprintln(c.outPutWriter, "Desc\t\t:", variable.AppDescLong)
	_, _ = fmt.Fprintln(c.outPutWriter, "Build Date\t:", version.BuildDate)
	_, _ = fmt.Fprintln(c.outPutWriter, "Git Commit\t:", version.CommitHash)
	_, _ = fmt.Fprintln(c.outPutWriter, "Version\t\t:", version.Version)
	_, _ = fmt.Fprintln(c.outPutWriter, "environment\t:", version.Environment)
	_, _ = fmt.Fprintln(c.outPutWriter, "Go Version\t:", variable.GoVersion)
	_, _ = fmt.Fprintln(c.outPutWriter, "OS / Arch\t:", variable.OsArch)

	return nil
}

func (c *CmdVersion) PersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	return nil
}

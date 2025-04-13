package process

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/helper"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/err"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/spf13/cobra"
)

type CmdProcess struct {
	outPutWriter io.Writer
	errWriter    io.Writer
	c            *cobra.Command
	wireApp      inject.Fn
	embedFS      *embed.FS
	appName      string
}

var _ cmd.Cmd = (*CmdProcess)(nil)

func NewCommand(appName string, wireApp inject.Fn, embedFS *embed.FS, outPutWriter io.Writer, errWriter io.Writer) *CmdProcess {
	return &CmdProcess{
		appName: appName,
		c: &cobra.Command{
			Use:     Usage,
			Short:   Short,
			Long:    Long,
			Aliases: Aliases,
			Example: fmt.Sprintf(
				"%s\n",
				fmt.Sprintf("Process data \n\t%s %s", appName, Example),
			),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		wireApp:      wireApp,
		embedFS:      embedFS,
		outPutWriter: outPutWriter,
		errWriter:    errWriter,
	}
}

func (c *CmdProcess) Init(_ *cmd.MetaData) *cobra.Command {
	c.c.PersistentPreRunE = c.PersistentPreRunner
	c.c.RunE = c.Runner

	c.c.SetOut(c.outPutWriter)
	c.c.SetErr(c.errWriter)
	c.initPersistentFlags()

	return c.c
}

func (c *CmdProcess) initPersistentFlags() {
	helper.InitCommonPersistentFlags(c.c)

	c.c.PersistentFlags().BoolVarP(
		&cmd.FlagIsDeleteCurrentReportDirectoryValue,
		cmd.FlagIsDeleteCurrentReportDirectory,
		cmd.FlagIsDeleteCurrentReportDirectoryShort,
		true,
		cmd.FlagIsDeleteCurrentReportDirectoryUsage,
	)

	c.c.PersistentFlags().StringVarP(
		&cmd.FlagReportTRXPathValue,
		cmd.FlagReportTRXPath,
		cmd.FlagReportTRXPathShort,
		filepath.Join(workDir, "report"),
		cmd.FlagReportTRXPathUsage,
	)
}

func (c *CmdProcess) Runner(_ *cobra.Command, _ []string) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	dBPath := csqlite.DBPath{}

	if cmd.FlagIsDebugValue {
		dBPath.WriteDBPath = "./reconciliation.db"
	}

	if app, cleanup, e := c.wireApp(
		c.c.Context(),
		c.embedFS,
		cconfig.AppName(c.appName),
		cconfig.TimeZone(cmd.FlagTZValue),
		err.RegisteredErrorType,
		clogger.IsShowLog(cmd.FlagIsVerboseValue),
		dBPath,
	); e == nil {
		atexit.Add(cleanup)
		e = helper.UpdateCommonConfigFromFlags(app)
		if e != nil {
			return e
		}

		app.GetComponents().Config.Data.Reconciliation.Action = c.c.Use
		app.GetComponents().Config.Data.Reconciliation.IsDeleteCurrentReportDirectory = cmd.FlagIsDeleteCurrentReportDirectoryValue
		app.GetComponents().Config.Data.Reconciliation.ReportTRXPath = cmd.FlagReportTRXPathValue
		app.Start()
	} else {
		return e
	}

	return nil
}

func (c *CmdProcess) PersistentPreRunner(cCmd *cobra.Command, args []string) (er error) {
	return helper.CommonPersistentPreRunner(cCmd, args)
}

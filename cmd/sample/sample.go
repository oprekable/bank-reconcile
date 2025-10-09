package sample

import (
	"embed"
	"fmt"
	"io"

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

type CmdSample struct {
	outPutWriter io.Writer
	errWriter    io.Writer
	c            *cobra.Command
	wireApp      inject.Fn
	embedFS      *embed.FS
	appName      string
}

var _ cmd.Cmd = (*CmdSample)(nil)

func NewCommand(appName string, wireApp inject.Fn, embedFS *embed.FS, outPutWriter io.Writer, errWriter io.Writer) *CmdSample {
	return &CmdSample{
		appName: appName,
		c: &cobra.Command{
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		wireApp:      wireApp,
		embedFS:      embedFS,
		outPutWriter: outPutWriter,
		errWriter:    errWriter,
	}
}

func (c *CmdSample) Init(_ *cmd.MetaData) *cobra.Command {
	c.c.Use = Usage
	c.c.Short = Short
	c.c.Long = Long
	c.c.Aliases = Aliases
	c.c.Example = fmt.Sprintf(
		"%s\n",
		fmt.Sprintf("Generate sample \n\t%s %s", c.appName, Example),
	)

	c.c.PersistentPreRunE = c.PersistentPreRunner
	c.c.RunE = c.Runner

	c.c.SetOut(c.outPutWriter)
	c.c.SetErr(c.errWriter)

	c.initPersistentFlags()

	return c.c
}

func (c *CmdSample) initPersistentFlags() {
	helper.InitCommonPersistentFlags(c.c)

	c.c.PersistentFlags().BoolVarP(
		&cmd.FlagIsDeleteCurrentSampleDirectoryValue,
		cmd.FlagIsDeleteCurrentSampleDirectory,
		cmd.FlagIsDeleteCurrentSampleDirectoryShort,
		true,
		cmd.FlagIsDeleteCurrentSampleDirectoryUsage,
	)

	c.c.PersistentFlags().Int64VarP(
		&cmd.FlagTotalDataSampleToGenerateValue,
		cmd.FlagTotalDataSampleToGenerate,
		cmd.FlagTotalDataSampleToGenerateShort,
		cmd.DefaultTotalDataSampleToGenerate,
		cmd.FlagTotalDataSampleToGenerateUsage,
	)

	c.c.PersistentFlags().IntVarP(
		&cmd.FlagPercentageMatchSampleToGenerateValue,
		cmd.FlagPercentageMatchSampleToGenerate,
		cmd.FlagPercentageMatchSampleToGenerateShort,
		cmd.DefaultPercentageMatchSampleToGenerate,
		cmd.FlagPercentageMatchSampleToGenerateUsage,
	)
}

func (c *CmdSample) Runner(_ *cobra.Command, _ []string) (er error) {
	dBPath := csqlite.DBPath{}

	if cmd.FlagIsDebugValue {
		dBPath.ReadDBPath = "./sample.db"
	}

	var app, cleanup, e = c.wireApp(
		c.c.Context(),
		c.embedFS,
		cconfig.AppName(c.appName),
		cconfig.TimeZone(cmd.FlagTZValue),
		err.RegisteredErrorType,
		clogger.IsShowLog(cmd.FlagIsVerboseValue),
		dBPath,
	)

	if e != nil {
		return e
	}

	atexit.Add(cleanup)

	e = helper.UpdateCommonConfigFromFlags(app)

	if e != nil {
		return e
	}

	conf := app.GetComponents().Config.Data
	conf.Reconciliation.Action = c.c.Use
	conf.Reconciliation.IsDeleteCurrentSampleDirectory = cmd.FlagIsDeleteCurrentSampleDirectoryValue
	conf.Reconciliation.TotalData = cmd.FlagTotalDataSampleToGenerateValue
	conf.Reconciliation.PercentageMatch = cmd.FlagPercentageMatchSampleToGenerateValue

	return app.Start()
}

func (c *CmdSample) PersistentPreRunner(cCmd *cobra.Command, args []string) (er error) {
	if cmd.FlagTotalDataSampleToGenerateValue <= 0 {
		return fmt.Errorf("'-%s' '--%s': %v should bigger than 0", cmd.FlagTotalDataSampleToGenerateShort, cmd.FlagTotalDataSampleToGenerate, cmd.FlagTotalDataSampleToGenerateValue)
	}

	if cmd.FlagPercentageMatchSampleToGenerateValue < 0 || cmd.FlagPercentageMatchSampleToGenerateValue > 100 {
		return fmt.Errorf("'-%s' '--%s': %v should between 0 and 100", cmd.FlagPercentageMatchSampleToGenerateShort, cmd.FlagPercentageMatchSampleToGenerate, cmd.FlagPercentageMatchSampleToGenerateValue)
	}

	return helper.CommonPersistentPreRunner(cCmd, args)
}

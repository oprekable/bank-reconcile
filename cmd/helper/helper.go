package helper

import (
	"context"
	"embed"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/err/core"
	"strconv"
	"time"

	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/spf13/cobra"
)

type IRunner interface {
	Run() (er error)
}

type WireApp func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error)
type Runner struct {
	wireApp WireApp
	cmd     *cobra.Command
	args    []string
}

func NewRunner(wireApp WireApp, cmd *cobra.Command, args []string) *Runner {
	return &Runner{
		wireApp: wireApp,
		cmd:     cmd,
		args:    args,
	}
}

func (r *Runner) Run(embedFs *embed.FS, appName string, tz string, errTypes []core.ErrorType, isVerbose bool, dbPath csqlite.DBPath) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	app, cleanup, er := r.wireApp(
		r.cmd.Context(),
		embedFs,
		cconfig.AppName(appName),
		cconfig.TimeZone(tz),
		errTypes,
		clogger.IsShowLog(isVerbose),
		dbPath,
	)

	if er != nil {
		return er
	}

	atexit.Add(cleanup)

	app.GetComponents().Config.Data.Reconciliation.Action = r.cmd.Use
	app.GetComponents().Config.Data.Reconciliation.SystemTRXPath = root.FlagSystemTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.BankTRXPath = root.FlagBankTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ReportTRXPath = root.FlagReportTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ListBank = root.FlagListBankValue
	app.GetComponents().Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory = root.FlagIsDeleteCurrentSampleDirectoryValue
	app.GetComponents().Config.Data.App.IsShowLog = strconv.FormatBool(root.FlagIsVerboseValue)
	app.GetComponents().Config.Data.App.IsDebug = root.FlagIsDebugValue
	app.GetComponents().Config.Data.App.IsProfilerActive = root.FlagIsProfilerActiveValue

	toDate, er := time.Parse(root.DateFormatString, root.FlagToDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.ToDate = toDate

	fromDate, er := time.Parse(root.DateFormatString, root.FlagFromDateValue)
	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.FromDate = fromDate
	app.GetComponents().Config.Data.Reconciliation.TotalData = root.FlagTotalDataSampleToGenerateValue
	app.GetComponents().Config.Data.Reconciliation.PercentageMatch = root.FlagPercentageMatchSampleToGenerateValue
	app.Start()

	return nil
}

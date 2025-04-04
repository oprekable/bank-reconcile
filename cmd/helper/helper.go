package helper

import (
	"strconv"
	"time"

	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/err"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

func RunnerSubCommand(cmd *cobra.Command, _ []string, dBPath csqlite.DBPath) (er error) {
	defer func() {
		atexit.AtExit()
	}()

	app, cleanup, er := inject.WireApp(
		cmd.Context(),
		root.EmbedFS,
		cconfig.AppName(variable.AppName),
		cconfig.TimeZone(root.FlagTZValue),
		err.RegisteredErrorType,
		clogger.IsShowLog(root.FlagIsVerboseValue),
		dBPath,
	)

	atexit.Add(cleanup)

	if er != nil {
		return er
	}

	app.GetComponents().Config.Data.Reconciliation.Action = cmd.Use
	app.GetComponents().Config.Data.Reconciliation.SystemTRXPath = root.FlagSystemTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.BankTRXPath = root.FlagBankTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ReportTRXPath = root.FlagReportTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ListBank = root.FlagListBankValue
	app.GetComponents().Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory = root.FlagIsDeleteCurrentSampleDirectoryValue
	app.GetComponents().Config.Data.IsShowLog = strconv.FormatBool(root.FlagIsVerboseValue)
	app.GetComponents().Config.Data.IsDebug = root.FlagIsDebugValue
	app.GetComponents().Config.Data.IsProfilerActive = root.FlagIsProfilerActiveValue

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

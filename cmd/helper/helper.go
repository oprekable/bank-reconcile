package helper

import (
	"fmt"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/oprekable/bank-reconcile/variable"
	"github.com/spf13/cobra"
	"path/filepath"
	"time"
)

func InitCommonPersistentFlags(c *cobra.Command) {
	defaultTZ := variable.TimeZone
	if defaultTZ == "" {
		defaultTZ = "Asia/Jakarta"
	}

	c.PersistentFlags().StringVarP(
		&cmd.FlagTZValue,
		cmd.FlagTimeZone,
		cmd.FlagTimeZoneShort,
		defaultTZ,
		cmd.FlagTimeZoneUsage,
	)

	workDir := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})

	c.PersistentFlags().StringVarP(
		&cmd.FlagSystemTRXPathValue,
		cmd.FlagSystemTRXPath,
		cmd.FlagSystemTRXPathShort,
		filepath.Join(workDir, "sample", "system"),
		cmd.FlagSystemTRXPathUsage,
	)

	c.PersistentFlags().StringVarP(
		&cmd.FlagBankTRXPathValue,
		cmd.FlagBankTRXPath,
		cmd.FlagBankTRXPathShort,
		filepath.Join(workDir, "sample", "bank"),
		cmd.FlagBankTRXPathUsage,
	)

	nowDateString := time.Now().Format("2006-01-02")

	c.PersistentFlags().StringVarP(
		&cmd.FlagFromDateValue,
		cmd.FlagFromDate,
		cmd.FlagFromDateShort,
		nowDateString,
		cmd.FlagFromDateUsage,
	)

	c.PersistentFlags().StringVarP(
		&cmd.FlagToDateValue,
		cmd.FlagToDate,
		cmd.FlagToDateShort,
		nowDateString,
		cmd.FlagToDateUsage,
	)

	c.PersistentFlags().StringSliceVarP(
		&cmd.FlagListBankValue,
		cmd.FlagListBank,
		cmd.FlagListBankShort,
		cmd.DefaultListBank,
		cmd.FlagListBankUsage,
	)

	c.PersistentFlags().BoolVarP(
		&cmd.FlagIsVerboseValue,
		cmd.FlagIsVerbose,
		cmd.FlagIsVerboseShort,
		false,
		cmd.FlagIsVerboseUsage,
	)

	c.PersistentFlags().BoolVarP(
		&cmd.FlagIsDebugValue,
		cmd.FlagIsDebug,
		cmd.FlagIsDebugShort,
		false,
		cmd.FlagIsDebugUsage,
	)

	c.PersistentFlags().BoolVarP(
		&cmd.FlagIsProfilerActiveValue,
		cmd.FlagIsProfilerActive,
		cmd.FlagIsProfilerActiveShort,
		false,
		cmd.FlagIsProfilerActiveUsage,
	)
}

func CommonPersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	if _, e := time.LoadLocation(cmd.FlagTZValue); e != nil {
		return fmt.Errorf("invalid value flag '-%s' '--%s': %v", cmd.FlagTimeZoneShort, cmd.FlagTimeZone, cmd.FlagTZValue)
	}

	fromDate, errFrom := time.Parse(cmd.DateFormatString, cmd.FlagFromDateValue)

	if errFrom != nil {
		return fmt.Errorf("failed to parse flag '-%s' '--%s': %v", cmd.FlagFromDateShort, cmd.FlagFromDate, cmd.FlagFromDateValue)
	}

	toDate, errTo := time.Parse(cmd.DateFormatString, cmd.FlagToDateValue)
	if errTo != nil {
		return fmt.Errorf("failed to parse flag '-%s' '--%s': %v", cmd.FlagToDateShort, cmd.FlagToDate, cmd.FlagToDateValue)
	}

	if fromDate.After(toDate) {
		return fmt.Errorf("'-%s' '--%s': %v should before '-%s' '--%s': %v", cmd.FlagFromDateShort, cmd.FlagFromDate, cmd.FlagFromDateValue, cmd.FlagToDateShort, cmd.FlagToDate, cmd.FlagToDateValue)
	}

	return nil
}

func UpdateCommonConfigFromFlags(app appcontext.IAppContext) (e error) {
	app.GetComponents().Config.Data.App.IsShowLog = cmd.FlagIsVerboseValue
	app.GetComponents().Config.Data.App.IsDebug = cmd.FlagIsDebugValue
	app.GetComponents().Config.Data.App.IsProfilerActive = cmd.FlagIsProfilerActiveValue

	app.GetComponents().Config.Data.Reconciliation.SystemTRXPath = cmd.FlagSystemTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.BankTRXPath = cmd.FlagBankTRXPathValue
	app.GetComponents().Config.Data.Reconciliation.ListBank = cmd.FlagListBankValue

	var toDate time.Time
	toDate, e = time.Parse(cmd.DateFormatString, cmd.FlagToDateValue)

	if e != nil {
		return e
	}

	app.GetComponents().Config.Data.Reconciliation.ToDate = toDate

	var fromDate time.Time
	fromDate, e = time.Parse(cmd.DateFormatString, cmd.FlagFromDateValue)

	if e != nil {
		return e
	}

	app.GetComponents().Config.Data.Reconciliation.FromDate = fromDate

	return nil
}

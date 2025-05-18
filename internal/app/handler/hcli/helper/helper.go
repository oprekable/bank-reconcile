package helper

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/schollz/progressbar/v3"
)

func InitProgressBar(writer io.Writer) *progressbar.ProgressBar {
	return progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(writer),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSpinnerType(17),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func InitCommonArgs(conf *config.Data, extraArgs [][]string) [][]string {
	formatText := "-%s --%s"
	args := [][]string{
		{
			fmt.Sprintf(formatText, cmd.FlagTimeZoneShort, cmd.FlagTimeZone),
			cmd.FlagTZValue,
		},
		{
			fmt.Sprintf(formatText, cmd.FlagFromDateShort, cmd.FlagFromDate),
			conf.Reconciliation.FromDate.Format("2006-01-02"),
		},
		{
			fmt.Sprintf(formatText, cmd.FlagToDateShort, cmd.FlagToDate),
			conf.Reconciliation.ToDate.Format("2006-01-02"),
		},
		{
			fmt.Sprintf(formatText, cmd.FlagSystemTRXPathShort, cmd.FlagSystemTRXPath),
			conf.Reconciliation.SystemTRXPath,
		},
		{
			fmt.Sprintf(formatText, cmd.FlagBankTRXPathShort, cmd.FlagBankTRXPath),
			conf.Reconciliation.BankTRXPath,
		},
		{
			fmt.Sprintf(formatText, cmd.FlagListBankShort, cmd.FlagListBank),
			strings.Join(conf.Reconciliation.ListBank, ","),
		},
		{
			fmt.Sprintf(formatText, cmd.FlagIsVerboseShort, cmd.FlagIsVerbose),
			strconv.FormatBool(conf.App.IsShowLog),
		},
		{
			fmt.Sprintf(formatText, cmd.FlagIsDebugShort, cmd.FlagIsDebug),
			strconv.FormatBool(conf.App.IsDebug),
		},
		{
			fmt.Sprintf(formatText, cmd.FlagIsProfilerActiveShort, cmd.FlagIsProfilerActive),
			strconv.FormatBool(conf.App.IsProfilerActive),
		},
	}

	args = append(args, extraArgs...)
	return args
}

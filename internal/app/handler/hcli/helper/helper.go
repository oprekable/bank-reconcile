package helper

import (
	"fmt"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"io"
	"strconv"
	"strings"

	"github.com/oprekable/bank-reconcile/cmd/root"

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
			fmt.Sprintf(formatText, root.FlagFromDateShort, root.FlagFromDate),
			conf.Reconciliation.FromDate.Format("2006-01-02"),
		},
		{
			fmt.Sprintf(formatText, root.FlagToDateShort, root.FlagToDate),
			conf.Reconciliation.ToDate.Format("2006-01-02"),
		},
		{
			fmt.Sprintf(formatText, root.FlagSystemTRXPathShort, root.FlagSystemTRXPath),
			conf.Reconciliation.SystemTRXPath,
		},
		{
			fmt.Sprintf(formatText, root.FlagBankTRXPathShort, root.FlagBankTRXPath),
			conf.Reconciliation.BankTRXPath,
		},
		{
			fmt.Sprintf(formatText, root.FlagListBankShort, root.FlagListBank),
			strings.Join(conf.Reconciliation.ListBank, ","),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsVerboseShort, root.FlagIsVerbose),
			strconv.FormatBool(conf.App.IsShowLog),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsDebugShort, root.FlagIsDebug),
			strconv.FormatBool(conf.App.IsDebug),
		},
		{
			fmt.Sprintf(formatText, root.FlagIsProfilerActiveShort, root.FlagIsProfilerActive),
			strconv.FormatBool(conf.App.IsProfilerActive),
		},
	}

	args = append(args, extraArgs...)
	return args
}

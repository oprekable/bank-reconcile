package cmd

import (
	"embed"
	"fmt"
	"path/filepath"
	"time"

	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   variable.AppName,
	Short: variable.AppDescShort,
	Long:  variable.AppDescLong,
	Example: fmt.Sprintf(
		"%s\n%s\n",
		fmt.Sprintf("Generate sample \n\t%s sample %s", variable.AppName, root.SampleUsageFlags),
		fmt.Sprintf("Process data \n\t%s process %s", variable.AppName, root.ProcessUsageFlags),
	),
	RunE:              root.Runner,
	PersistentPreRunE: root.PersistentPreRunner,
}

func init() {
	defaultTZ := variable.TimeZone
	if defaultTZ == "" {
		defaultTZ = "Asia/Jakarta"
	}

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagTZValue,
		root.FlagTimeZone,
		root.FlagTimeZoneShort,
		defaultTZ,
		root.FlagTimeZoneUsage,
	)

	workDir := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
	rootCmd.PersistentFlags().StringVarP(
		&root.FlagSystemTRXPathValue,
		root.FlagSystemTRXPath,
		root.FlagSystemTRXPathShort,
		filepath.Join(workDir, "sample", "system"),
		root.FlagSystemTRXPathUsage,
	)

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagBankTRXPathValue,
		root.FlagBankTRXPath,
		root.FlagBankTRXPathShort,
		filepath.Join(workDir, "sample", "bank"),
		root.FlagBankTRXPathUsage,
	)

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagReportTRXPathValue,
		root.FlagReportTRXPath,
		root.FlagReportTRXPathShort,
		filepath.Join(workDir, "report"),
		root.FlagReportTRXPathUsage,
	)

	nowDateString := time.Now().Format("2006-01-02")

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagFromDateValue,
		root.FlagFromDate,
		root.FlagFromDateShort,
		nowDateString,
		root.FlagFromDateUsage,
	)

	rootCmd.PersistentFlags().StringVarP(
		&root.FlagToDateValue,
		root.FlagToDate,
		root.FlagToDateShort,
		nowDateString,
		root.FlagToDateUsage,
	)

	rootCmd.PersistentFlags().StringSliceVarP(
		&root.FlagListBankValue,
		root.FlagListBank,
		root.FlagListBankShort,
		root.DefaultListBank,
		root.FlagListBankUsage,
	)

	rootCmd.PersistentFlags().Int64VarP(
		&root.FlagTotalDataSampleToGenerateValue,
		root.FlagTotalDataSampleToGenerate,
		root.FlagTotalDataSampleToGenerateShort,
		root.DefaultTotalDataSampleToGenerate,
		root.FlagTotalDataSampleToGenerateUsage,
	)

	rootCmd.PersistentFlags().IntVarP(
		&root.FlagPercentageMatchSampleToGenerateValue,
		root.FlagPercentageMatchSampleToGenerate,
		root.FlagPercentageMatchSampleToGenerateShort,
		root.DefaultPercentageMatchSampleToGenerate,
		root.FlagPercentageMatchSampleToGenerateUsage,
	)

	rootCmd.PersistentFlags().BoolVarP(
		&root.FlagIsDeleteCurrentSampleDirectoryValue,
		root.FlagIsDeleteCurrentSampleDirectory,
		root.FlagIsDeleteCurrentSampleDirectoryShort,
		true,
		root.FlagIsDeleteCurrentSampleDirectoryUsage,
	)

	rootCmd.PersistentFlags().BoolVarP(
		&root.FlagIsVerboseValue,
		root.FlagIsVerbose,
		root.FlagIsVerboseShort,
		false,
		root.FlagIsVerboseUsage,
	)

	rootCmd.PersistentFlags().BoolVarP(
		&root.FlagIsDebugValue,
		root.FlagIsDebug,
		root.FlagIsDebugShort,
		false,
		root.FlagIsDebugUsage,
	)

	rootCmd.PersistentFlags().BoolVarP(
		&root.FlagIsProfilerActiveValue,
		root.FlagIsProfilerActive,
		root.FlagIsProfilerActiveShort,
		false,
		root.FlagIsProfilerActiveUsage,
	)
}

func Execute(embedFS *embed.FS) {
	root.EmbedFS = embedFS
	if er := rootCmd.Execute(); er != nil {
		fmt.Println(er)
	}
}

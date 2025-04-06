package root

import (
	"embed"
	"fmt"
	"path/filepath"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
)

const (
	DateFormatString                         string = "2006-01-02"
	FlagSystemTRXPath                        string = "systemtrxpath"
	FlagSystemTRXPathShort                   string = "s"
	FlagSystemTRXPathUsage                   string = "Path location of System Transaction directory"
	FlagBankTRXPath                          string = "banktrxpath"
	FlagBankTRXPathShort                     string = "b"
	FlagBankTRXPathUsage                     string = "Path location of Bank Transaction directory"
	FlagReportTRXPath                        string = "reportpath"
	FlagReportTRXPathShort                   string = "r"
	FlagReportTRXPathUsage                   string = "Path location of Archive directory"
	FlagListBank                             string = "listbank"
	FlagListBankShort                        string = "l"
	FlagListBankUsage                        string = "List bank accepted"
	FlagTimeZone                             string = "time_zone"
	FlagTimeZoneShort                        string = "z"
	FlagTimeZoneUsage                        string = `time zone settings`
	FlagFromDate                             string = "from"
	FlagFromDateShort                        string = "f"
	FlagFromDateUsage                        string = `from date (YYYY-MM-DD)`
	FlagToDate                               string = "to"
	FlagToDateShort                          string = "t"
	FlagToDateUsage                          string = `to date (YYYY-MM-DD)`
	FlagTotalDataSampleToGenerate            string = "amountdata"
	FlagTotalDataSampleToGenerateShort       string = "a"
	FlagTotalDataSampleToGenerateUsage       string = `amount system trx data sample to generate, bank trx will be 2 times of this amount`
	FlagPercentageMatchSampleToGenerate      string = "percentagematch"
	FlagPercentageMatchSampleToGenerateShort string = "p"
	FlagPercentageMatchSampleToGenerateUsage string = `percentage of matched trx for data sample to generate`
	FlagIsDeleteCurrentSampleDirectory       string = "deleteoldfile"
	FlagIsDeleteCurrentSampleDirectoryShort  string = "d"
	FlagIsDeleteCurrentSampleDirectoryUsage  string = `delete old files`
	FlagIsVerbose                            string = "showlog"
	FlagIsVerboseShort                       string = "o"
	FlagIsVerboseUsage                       string = `show logs`
	FlagIsDebug                              string = "debug"
	FlagIsDebugShort                         string = "g"
	FlagIsDebugUsage                         string = `debug mode`
	FlagIsProfilerActive                     string = "profiler"
	FlagIsProfilerActiveShort                string = "i"
	FlagIsProfilerActiveUsage                string = `pprof active mode`
)

var EmbedFS *embed.FS
var nowDateString = time.Now().Format("2006-01-02")
var workDir = filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
var pathSystemTrx = filepath.Join(workDir, "sample", "system")
var pathBankTrx = filepath.Join(workDir, "sample", "bank")
var pathReportTrx = filepath.Join(workDir, "report")
var SampleUsageFlags = fmt.Sprintf("--systemtrxpath=%s --banktrxpath=%s --listbank=bca,bni,mandiri,bri,danamon --percentagematch=100 --amountdata=10000 --from=%s --to=%s", pathSystemTrx, pathBankTrx, nowDateString, nowDateString)
var ProcessUsageFlags = fmt.Sprintf("--systemtrxpath=%s --banktrxpath=%s --reportpath=%s --listbank=bca,bni,mandiri,bri,danamon --from=%s --to=%s", pathSystemTrx, pathBankTrx, pathReportTrx, nowDateString, nowDateString)
var FlagSystemTRXPathValue string
var FlagBankTRXPathValue string
var FlagReportTRXPathValue string
var FlagListBankValue []string
var FlagTZValue string
var FlagFromDateValue string
var FlagToDateValue string
var DefaultListBank = []string{"bca", "bni", "mandiri", "bri", "danamon"}
var FlagTotalDataSampleToGenerateValue int64
var DefaultTotalDataSampleToGenerate int64 = 1000
var FlagPercentageMatchSampleToGenerateValue int
var DefaultPercentageMatchSampleToGenerate = 100
var FlagIsDeleteCurrentSampleDirectoryValue bool
var FlagIsVerboseValue bool
var FlagIsDebugValue bool
var FlagIsProfilerActiveValue bool

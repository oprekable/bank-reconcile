package cmd

var FlagTZValue string
var FlagBankTRXPathValue string
var FlagSystemTRXPathValue string
var FlagFromDateValue string
var FlagToDateValue string
var FlagListBankValue []string
var DefaultListBank = []string{"bca", "bni", "mandiri", "bri", "danamon"}
var FlagTotalDataSampleToGenerateValue int64
var DefaultTotalDataSampleToGenerate int64 = 1000
var FlagPercentageMatchSampleToGenerateValue int
var DefaultPercentageMatchSampleToGenerate = 100
var FlagIsDeleteCurrentSampleDirectoryValue bool
var FlagIsDeleteCurrentReportDirectoryValue bool
var FlagIsVerboseValue bool
var FlagIsDebugValue bool
var FlagIsProfilerActiveValue bool
var FlagReportTRXPathValue string

const (
	DateFormatString                         string = "2006-01-02"
	FlagTimeZone                             string = "time_zone"
	FlagTimeZoneShort                        string = "z"
	FlagTimeZoneUsage                        string = `time zone settings`
	FlagSystemTRXPath                        string = "systemtrxpath"
	FlagSystemTRXPathShort                   string = "s"
	FlagSystemTRXPathUsage                   string = "Path location of System Transaction directory"
	FlagBankTRXPath                          string = "banktrxpath"
	FlagBankTRXPathShort                     string = "b"
	FlagBankTRXPathUsage                     string = "Path location of Bank Transaction directory"
	FlagReportTRXPath                        string = "reportpath"
	FlagReportTRXPathShort                   string = "r"
	FlagReportTRXPathUsage                   string = "Path location of Archive directory"
	FlagFromDate                             string = "from"
	FlagFromDateShort                        string = "f"
	FlagFromDateUsage                        string = `from date (YYYY-MM-DD)`
	FlagToDate                               string = "to"
	FlagToDateShort                          string = "t"
	FlagToDateUsage                          string = `to date (YYYY-MM-DD)`
	FlagListBank                             string = "listbank"
	FlagListBankShort                        string = "l"
	FlagListBankUsage                        string = "List bank accepted"
	FlagTotalDataSampleToGenerate            string = "amountdata"
	FlagTotalDataSampleToGenerateShort       string = "a"
	FlagTotalDataSampleToGenerateUsage       string = `amount system trx data sample to generate, bank trx will be 2 times of this amount`
	FlagPercentageMatchSampleToGenerate      string = "percentagematch"
	FlagPercentageMatchSampleToGenerateShort string = "p"
	FlagPercentageMatchSampleToGenerateUsage string = `percentage of matched trx for data sample to generate`
	FlagIsDeleteCurrentSampleDirectory       string = "deleteoldfile"
	FlagIsDeleteCurrentSampleDirectoryShort  string = "d"
	FlagIsDeleteCurrentSampleDirectoryUsage  string = `delete old sample files`
	FlagIsDeleteCurrentReportDirectory       string = "deleteoldfile"
	FlagIsDeleteCurrentReportDirectoryShort  string = "d"
	FlagIsDeleteCurrentReportDirectoryUsage  string = `delete old report files`
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

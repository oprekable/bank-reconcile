package process

type FilePathBankTrx struct {
	Bank     string
	FilePath string
}

type ReconciliationSummary struct {
	FileMissingBankTrx              map[string]string `deepcopier:"skip"`
	FileMissingSystemTrx            string            `deepcopier:"skip"`
	FileMatchedSystemTrx            string            `deepcopier:"skip"`
	TotalProcessedSystemTrx         int64             `deepcopier:"field:TotalSystemTrx"`
	TotalMatchedSystemTrx           int64             `deepcopier:"field:TotalMatchedTrx"`
	TotalNotMatchedSystemTrx        int64             `deepcopier:"field:TotalNotMatchedTrx"`
	SumAmountProcessedSystemTrx     float64           `deepcopier:"field:SumSystemTrx"`
	SumAmountMatchedSystemTrx       float64           `deepcopier:"field:SumMatchedTrx"`
	SumAmountDiscrepanciesSystemTrx float64           `deepcopier:"field:SumDiscrepanciesTrx"`
}

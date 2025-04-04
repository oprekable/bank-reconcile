package process

type ReconciliationSummary struct {
	TotalSystemTrx      int64   `db:"total_system_trx"`
	TotalMatchedTrx     int64   `db:"total_matched_trx"`
	TotalNotMatchedTrx  int64   `db:"total_not_matched_trx"`
	SumSystemTrx        float64 `db:"sum_system_trx"`
	SumMatchedTrx       float64 `db:"sum_matched_trx"`
	SumDiscrepanciesTrx float64 `db:"sum_discrepancies_trx"`
}

type MatchedTrx struct {
	SystemTrxTrxID           string  `db:"SystemTrxTrxID"`
	BankTrxUniqueIdentifier  string  `db:"BankTrxUniqueIdentifier"`
	SystemTrxTransactionTime string  `db:"SystemTrxTransactionTime"`
	BankTrxDate              string  `db:"BankTrxDate"`
	SystemTrxType            string  `db:"SystemTrxType"`
	Bank                     string  `db:"Bank"`
	SystemTrxAmount          float64 `db:"SystemTrxAmount"`
	BankTrxAmount            float64 `db:"BankTrxAmount"`
}

type NotMatchedSystemTrx struct {
	TrxID           string  `db:"TrxID"`
	TransactionTime string  `db:"TransactionTime"`
	Type            string  `db:"Type"`
	Amount          float64 `db:"Amount"`
}

type NotMatchedBankTrx struct {
	UniqueIdentifier string  `db:"UniqueIdentifier"`
	Bank             string  `db:"Bank"`
	Date             string  `db:"Date"`
	Amount           float64 `db:"Amount"`
}

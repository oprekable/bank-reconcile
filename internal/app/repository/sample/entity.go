package sample

type TrxData struct {
	TrxID            string  `db:"trxID"`
	UniqueIdentifier string  `db:"uniqueIdentifier"`
	Type             string  `db:"type"`
	Bank             string  `db:"bank"`
	TransactionTime  string  `db:"transactionTime"`
	Date             string  `db:"date"`
	IsSystemTrx      bool    `db:"is_system_trx"`
	IsBankTrx        bool    `db:"is_bank_trx"`
	Amount           float64 `db:"amount"`
}

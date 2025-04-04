package banks

import "time"

type BankParserType string

const (
	DefaultBankParser BankParserType = "DEFAULT"
	BCABankParser     BankParserType = "BCA"
	BNIBankParser     BankParserType = "BNI"
)

type TrxType string

const (
	DEBIT  TrxType = "DEBIT"
	CREDIT TrxType = "CREDIT"
)

type BankTrxData struct {
	UniqueIdentifier string
	Date             time.Time
	Type             TrxType
	Bank             string
	FilePath         string
	Amount           float64
}

package banks

import "context"

type ReconcileBankData interface {
	GetBank() string
	GetParser() BankParserType
	ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error)
}

type BankTrxDataInterface interface {
	GetUniqueIdentifier() string
	GetDate() string
	GetAmount() float64
	GetAbsAmount() float64
	GetType() TrxType
	GetBank() string
	ToBankTrxData() (returnData *BankTrxData, err error)
}

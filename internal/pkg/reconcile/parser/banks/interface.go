package banks

import "context"

//go:generate mockery --name "ReconcileBankData" --output "./_mock" --outpkg "_mock"
type ReconcileBankData interface {
	GetBank() string
	GetParser() BankParserType
	ToBankTrxData(ctx context.Context, filePath string) (returnData []*BankTrxData, err error)
}

//go:generate mockery --name "BankTrxDataInterface" --output "./_mock" --outpkg "_mock"
type BankTrxDataInterface interface {
	GetUniqueIdentifier() string
	GetDate() string
	GetAmount() float64
	GetAbsAmount() float64
	GetType() TrxType
	GetBank() string
	ToBankTrxData() (returnData *BankTrxData, err error)
}

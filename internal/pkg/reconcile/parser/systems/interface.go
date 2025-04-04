package systems

import "context"

//go:generate mockery --name "ReconcileSystemData" --output "./_mock" --outpkg "_mock"
type ReconcileSystemData interface {
	ToSystemTrxData(ctx context.Context, filePath string) (returnData []*SystemTrxData, err error)
}

//go:generate mockery --name "SystemTrxDataInterface" --output "./_mock" --outpkg "_mock"
type SystemTrxDataInterface interface {
	GetTrxID() string
	GetTransactionTime() string
	GetAmount() float64
	GetType() TrxType
	ToSystemTrxData() (returnData *SystemTrxData, err error)
}

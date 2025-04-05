package systems

import "context"

type ReconcileSystemData interface {
	ToSystemTrxData(ctx context.Context, filePath string) (returnData []*SystemTrxData, err error)
}

type SystemTrxDataInterface interface {
	GetTrxID() string
	GetTransactionTime() string
	GetAmount() float64
	GetType() TrxType
	ToSystemTrxData() (returnData *SystemTrxData, err error)
}

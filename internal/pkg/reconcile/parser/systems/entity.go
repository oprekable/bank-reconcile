package systems

import "time"

type SystemParserType string

const (
	DefaultSystemParser SystemParserType = "DEFAULT"
)

type TrxType string

type SystemTrxData struct {
	TrxID           string
	TransactionTime time.Time
	Type            TrxType
	FilePath        string
	Amount          float64
}

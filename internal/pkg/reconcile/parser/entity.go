package parser

import (
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
)

type TrxData struct {
	SystemTrx       []*systems.SystemTrxData
	BankTrx         []*banks.BankTrxData
	MinSystemAmount float64
	MaxSystemAmount float64
}

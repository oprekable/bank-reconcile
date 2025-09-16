package service

import (
	"encoding/csv"

	"github.com/google/wire"
	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	"github.com/oprekable/bank-reconcile/internal/app/service/sample"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bca"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bni"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/default_bank"
)

// ProvideBankParserFactoryMap creates the map of all available bank parser factories.
// This is the correct, high-level location for assembling all parser implementations.
func ProvideBankParserFactoryMap() map[string]banks.BankParserFactory {
	factories := make(map[string]banks.BankParserFactory)

	// Register BCA parser
	factories[string(banks.BCABankParser)] = func(bankName string, reader *csv.Reader, hasHeader bool) (banks.ReconcileBankData, error) {
		return bca.NewBankParser(bankName, reader, hasHeader)
	}

	// Register BNI parser
	factories[string(banks.BNIBankParser)] = func(bankName string, reader *csv.Reader, hasHeader bool) (banks.ReconcileBankData, error) {
		return bni.NewBankParser(bankName, reader, hasHeader)
	}

	// Register Default parser
	factories[string(banks.DefaultBankParser)] = func(bankName string, reader *csv.Reader, hasHeader bool) (banks.ReconcileBankData, error) {
		return default_bank.NewBankParser(bankName, reader, hasHeader)
	}

	return factories
}

func NewServices(
	svcSample sample.Service,
	svcProcess process.Service,
) *Services {
	return &Services{
		SvcSample:  svcSample,
		SvcProcess: svcProcess,
	}
}

var Set = wire.NewSet(
	// Provide the parser registry dependencies
	ProvideBankParserFactoryMap,
	banks.NewParserRegistry,

	// Provide the services
	sample.Set,
	process.Set,
	NewServices,
)

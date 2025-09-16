package banks

import (
	"encoding/csv"
	"fmt"
	"io"
)

// BankParserFactory defines the signature for a function that creates a new bank parser.
type BankParserFactory func(bankName string, reader *csv.Reader, hasHeader bool) (ReconcileBankData, error)

// ParserRegistry holds the collection of available bank parser factories.
// It is managed by the dependency injection container.
type ParserRegistry struct {
	factories map[string]BankParserFactory
}

// NewParserRegistry creates a new instance of ParserRegistry.
func NewParserRegistry(factories map[string]BankParserFactory) *ParserRegistry {
	return &ParserRegistry{factories: factories}
}

// GetParser retrieves a parser instance from the registry.
func (r *ParserRegistry) GetParser(bankName string, fileReader io.Reader, hasHeader bool) (ReconcileBankData, error) {
	factory, ok := r.factories[bankName]
	if !ok {
		// Fallback to a default parser if the specific one is not found
		defaultFactory, defaultOk := r.factories["DEFAULT"]
		if !defaultOk {
			return nil, fmt.Errorf("bank parser for '%s' not found and no default parser is registered", bankName)
		}
		return defaultFactory(bankName, csv.NewReader(fileReader), hasHeader)
	}
	return factory(bankName, csv.NewReader(fileReader), hasHeader)
}

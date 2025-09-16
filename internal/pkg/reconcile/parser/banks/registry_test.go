package banks

import (
	"context"
	"encoding/csv"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Mock Implementations for Testing ---

type mockParser struct {
	bankName   string
	parserType BankParserType
}

func (m *mockParser) ToBankTrxData(_ context.Context, _ string) ([]*BankTrxData, error) {
	return nil, errors.New("not implemented for mock")
}

func (m *mockParser) GetParser() BankParserType {
	return m.parserType
}

func (m *mockParser) GetBank() string {
	return m.bankName
}

// --- Test Cases ---

// TestParserRegistry tests the functionality of the ParserRegistry in isolation.
func TestParserRegistry(t *testing.T) {
	// 1. Setup: Create mock factories locally for the test.
	factories := make(map[string]BankParserFactory)

	factories["MOCK_BCA"] = func(bankName string, reader *csv.Reader, hasHeader bool) (ReconcileBankData, error) {
		return &mockParser{bankName: bankName, parserType: "MOCK_BCA_PARSER"}, nil
	}
	factories["DEFAULT"] = func(bankName string, reader *csv.Reader, hasHeader bool) (ReconcileBankData, error) {
		return &mockParser{bankName: bankName, parserType: DefaultBankParser}, nil
	}

	registry := NewParserRegistry(factories)

	// Dummy reader for tests, as the content doesn't matter for this test.
	dummyReader := strings.NewReader("")

	t.Run("should get a specific parser for a registered bank", func(t *testing.T) {
		parser, err := registry.GetParser("MOCK_BCA", dummyReader, true)
		assert.NoError(t, err)
		assert.NotNil(t, parser)
		assert.Equal(t, "MOCK_BCA", parser.GetBank())
		assert.Equal(t, BankParserType("MOCK_BCA_PARSER"), parser.GetParser())
	})

	t.Run("should fall back to default parser for an unregistered bank", func(t *testing.T) {
		parser, err := registry.GetParser("UNKNOWN_BANK", dummyReader, true)
		assert.NoError(t, err)
		assert.NotNil(t, parser)
		assert.Equal(t, "UNKNOWN_BANK", parser.GetBank())
		assert.Equal(t, DefaultBankParser, parser.GetParser())
	})

	t.Run("should return an error if no parser is found and no default is registered", func(t *testing.T) {
		// Setup for this specific case: Create a registry without a default parser.
		emptyFactories := make(map[string]BankParserFactory)
		emptyRegistry := NewParserRegistry(emptyFactories)

		parser, err := emptyRegistry.GetParser("ANYBANK", dummyReader, true)
		assert.Error(t, err)
		assert.Nil(t, parser)
		assert.Contains(t, err.Error(), "not found and no default parser is registered")
	})
}

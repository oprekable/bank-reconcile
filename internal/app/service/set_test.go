package service

import (
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/service/process/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/service/sample"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/service/sample/_mock"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestNewServices(t *testing.T) {
	type args struct {
		svcSample  sample.ServiceGenerator
		svcProcess process.ServiceGenerator
	}

	tests := []struct {
		args args
		want *Services
		name string
	}{
		{
			name: "Ok",
			args: args{
				svcSample:  mocksample.NewServiceGenerator(t),
				svcProcess: mockprocess.NewServiceGenerator(t),
			},
			want: &Services{
				SvcSample:  mocksample.NewServiceGenerator(t),
				SvcProcess: mockprocess.NewServiceGenerator(t),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServices(tt.args.svcSample, tt.args.svcProcess); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvideBankParserFactoryMap(t *testing.T) {
	tests := []struct {
		name       string
		bank       string
		parser     string
		wantParser string
	}{
		{
			name:       "DEFAULT ok",
			bank:       string(banks.DefaultBankParser),
			parser:     string(banks.DefaultBankParser),
			wantParser: string(banks.DefaultBankParser),
		},
		{
			name:       "BCA ok",
			bank:       string(banks.BCABankParser),
			parser:     string(banks.BCABankParser),
			wantParser: string(banks.BCABankParser),
		},
		{
			name:       "BNI ok",
			bank:       string(banks.BNIBankParser),
			parser:     string(banks.BNIBankParser),
			wantParser: string(banks.BNIBankParser),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProvideBankParserFactoryMap()
			parserFactory, ok := got[tt.parser]

			if !ok {
				t.Errorf("ProvideBankParserFactoryMap() %v not found", tt.parser)
			}

			reconcileBankData, _ := parserFactory(tt.bank, csv.NewReader(nil), true)

			if gotParser := reconcileBankData.GetParser(); string(gotParser) != tt.wantParser {
				t.Errorf("ProvideBankParserFactoryMap() = %v, want %v", gotParser, tt.wantParser)
			}
		})
	}
}

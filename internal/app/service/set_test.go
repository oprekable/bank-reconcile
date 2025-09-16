package service

import (
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
		svcSample  sample.Service
		svcProcess process.Service
	}

	tests := []struct {
		args args
		want *Services
		name string
	}{
		{
			name: "Ok",
			args: args{
				svcSample:  mocksample.NewService(t),
				svcProcess: mockprocess.NewService(t),
			},
			want: &Services{
				SvcSample:  mocksample.NewService(t),
				SvcProcess: mockprocess.NewService(t),
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
		name string
		want map[string]banks.BankParserFactory
	}{
		{
			name: "Ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProvideBankParserFactoryMap()

			if _, ok := got[string(banks.DefaultBankParser)]; !ok {
				t.Errorf("ProvideBankParserFactoryMap() DEFAULT not found")
			}

			if _, ok := got[string(banks.BCABankParser)]; !ok {
				t.Errorf("ProvideBankParserFactoryMap() BCA not found")
			}

			if _, ok := got[string(banks.BNIBankParser)]; !ok {
				t.Errorf("ProvideBankParserFactoryMap() BNI not found")
			}
		})
	}
}

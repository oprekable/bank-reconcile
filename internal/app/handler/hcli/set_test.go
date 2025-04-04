package hcli

import (
	"os"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/process"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/sample"
)

func TestProviderHandlers(t *testing.T) {
	tests := []struct {
		name string
		want []Handler
	}{
		{
			name: "Ok",
			want: []Handler{
				noop.NewHandler(os.Stdout),
				process.NewHandler(os.Stdout),
				sample.NewHandler(os.Stdout),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderHandlers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

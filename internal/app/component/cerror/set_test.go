package cerror

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/err/core"
)

func TestProvideErType(t *testing.T) {
	type args struct {
		errType []core.ErrorType
	}

	tests := []struct {
		name string
		args args
		want ErType
	}{
		{
			name: "Ok",
			args: args{
				errType: []core.ErrorType{
					core.CErrInternal,
					core.CErrDBConn,
				},
			},
			want: []core.ErrorType{
				core.CErrInternal,
				core.CErrDBConn,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProvideErType(tt.args.errType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProvideErType() = %v, want %v", got, tt.want)
			}
		})
	}
}

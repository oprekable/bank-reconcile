package cerror

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/err/core"
)

func TestErrorGetErrors(t *testing.T) {
	type fields struct {
		errors []error
	}

	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "Ok",
			fields: fields{
				errors: []error{
					core.CErrInternal.Error(),
					core.CErrDBConn.Error(),
				},
			},
			want: []error{
				core.CErrInternal.Error(),
				core.CErrDBConn.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				errors: tt.fields.errors,
			}

			if got := e.GetErrors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		erType ErType
	}

	tests := []struct {
		want *Error
		name string
		args args
	}{
		{
			name: "Ok",
			args: args{
				erType: ErType{
					core.CErrInternal,
					core.CErrDBConn,
				},
			},
			want: &Error{
				errors: []error{
					core.CErrInternal.Error(),
					core.CErrDBConn.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.erType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

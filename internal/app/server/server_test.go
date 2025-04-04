package server

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/server/_mock"

	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
)

func TestNewServer(t *testing.T) {
	type args struct {
		cli CliServer
	}

	tests := []struct {
		args args
		want *Server
		name string
	}{
		{
			name: "Ok",
			args: args{
				cli: nil,
			},
			want: NewServer(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.cli)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerRun(t *testing.T) {
	type fields struct {
		Cli CliServer
	}

	type args struct {
		eg *errgroup.Group
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				Cli: func() CliServer {
					m := _mock.NewIServer(t)
					m.On(
						"Start",
						mock.Anything,
					).Maybe()
					return m
				}(),
			},
			args: args{
				eg: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Cli: tt.fields.Cli,
			}

			s.Run(tt.args.eg)
		})
	}
}

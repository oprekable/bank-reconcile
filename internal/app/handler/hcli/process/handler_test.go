package process

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/service/process/_mock"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
)

func TestHandlerExec(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp   *component.Components
		svc    *service.Services
		repo   *repository.Repositories
		writer io.Writer
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "Nil components services repository",
			fields: fields{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
			wantErr: false,
		},
		{
			name: "Error GenerateReconciliation",
			fields: fields{
				comp: &component.Components{
					Context: context.TODO(),
					Fs: &cfs.Fs{
						LocalStorageFs: afero.NewMemMapFs(),
					},
					Config: &cconfig.Config{
						Data: &config.Data{
							App: core.App{},
							Reconciliation: reconciliation.Reconciliation{
								ReportTRXPath: "/foo",
							},
						},
					},
				},
				svc: func() *service.Services {
					mockSvc := mockprocess.NewServiceGenerator(t)

					mockSvc.On(
						"GenerateReconciliation",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(
						process.ReconciliationSummary{},
						errors.New("error"),
					).Maybe()

					return service.NewServices(
						nil,
						mockSvc,
					)
				}(),
				repo:   &repository.Repositories{},
				writer: &bf,
			},
			wantErr: true,
		},
		{
			name: "Ok",
			fields: fields{
				comp: &component.Components{
					Context: context.TODO(),
					Fs: &cfs.Fs{
						LocalStorageFs: afero.NewMemMapFs(),
					},
					Config: &cconfig.Config{
						Data: &config.Data{
							App: core.App{},
							Reconciliation: reconciliation.Reconciliation{
								ReportTRXPath: "/foo",
							},
						},
					},
				},
				svc: func() *service.Services {
					mockSvc := mockprocess.NewServiceGenerator(t)

					mockSvc.On(
						"GenerateReconciliation",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(
						process.ReconciliationSummary{
							FileMissingSystemTrx: "/foo.csv",
							FileMissingBankTrx: func() map[string]string {
								m := make(map[string]string)
								m["foo"] = "/bar.csv"
								return m
							}(),
						},
						nil,
					).Maybe()

					return service.NewServices(
						nil,
						mockSvc,
					)
				}(),
				repo:   &repository.Repositories{},
				writer: &bf,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: tt.fields.writer,
			}

			if err := h.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
			}

			bf.Reset()
		})
	}
}

func TestHandlerName(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp   *component.Components
		svc    *service.Services
		repo   *repository.Repositories
		writer io.Writer
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
			want: "process",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: tt.fields.writer,
			}

			if got := h.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestHandlerSetComponents(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp   *component.Components
		svc    *service.Services
		repo   *repository.Repositories
		writer io.Writer
	}

	type args struct {
		c *component.Components
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
			args: args{
				c: &component.Components{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: tt.fields.writer,
			}

			h.SetComponents(tt.args.c)

			bf.Reset()
		})
	}
}

func TestHandlerSetRepositories(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp   *component.Components
		svc    *service.Services
		repo   *repository.Repositories
		writer io.Writer
	}

	type args struct {
		r *repository.Repositories
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
			args: args{
				r: &repository.Repositories{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: tt.fields.writer,
			}

			h.SetRepositories(tt.args.r)

			bf.Reset()
		})
	}
}

func TestHandlerSetServices(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp   *component.Components
		svc    *service.Services
		repo   *repository.Repositories
		writer io.Writer
	}

	type args struct {
		s *service.Services
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
			args: args{
				s: &service.Services{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: tt.fields.writer,
			}

			h.SetServices(tt.args.s)

			bf.Reset()
		})
	}
}

func TestNewHandler(t *testing.T) {
	var bf bytes.Buffer
	tests := []struct {
		want *Handler
		name string
	}{
		{
			name: "Ok",
			want: &Handler{
				comp:   nil,
				svc:    nil,
				repo:   nil,
				writer: &bf,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(&bf)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

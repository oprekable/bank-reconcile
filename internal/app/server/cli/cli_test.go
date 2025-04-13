package cli

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
)

func TestCliName(t *testing.T) {
	type fields struct {
		ctx      context.Context
		comp     *component.Components
		svc      *service.Services
		repo     *repository.Repositories
		handlers []hcli.Handler
	}

	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:      nil,
				comp:     nil,
				svc:      nil,
				repo:     nil,
				handlers: nil,
			},
			want: "cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{
				ctx:      tt.fields.ctx,
				comp:     tt.fields.comp,
				svc:      tt.fields.svc,
				repo:     tt.fields.repo,
				handlers: tt.fields.handlers,
			}

			if got := c.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCliShutdown(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		ctx      context.Context
		comp     *component.Components
		svc      *service.Services
		repo     *repository.Repositories
		handlers []hcli.Handler
	}

	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:      zerolog.New(&bf).WithContext(context.Background()),
				comp:     nil,
				svc:      nil,
				repo:     nil,
				handlers: nil,
			},
			want: `{"level":"info","message":"[cli] shutdown"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{
				ctx:      tt.fields.ctx,
				comp:     tt.fields.comp,
				svc:      tt.fields.svc,
				repo:     tt.fields.repo,
				handlers: tt.fields.handlers,
			}

			c.Shutdown()

			if got := bf.String(); strings.TrimRight(got, "\n") != tt.want {
				t.Errorf("Msg() output = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCliStart(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		ctx      context.Context
		comp     *component.Components
		svc      *service.Services
		repo     *repository.Repositories
		handlers []hcli.Handler
	}

	type args struct {
		eg *errgroup.Group
	}

	tests := []struct {
		args   args
		name   string
		fields fields
	}{
		{
			name: "Ok",
			fields: fields{
				ctx: context.Background(),
				comp: &component.Components{
					Config: &cconfig.Config{
						Data: &config.Data{
							Reconciliation: reconciliation.Reconciliation{
								Action: "noop",
							},
						},
					},
				},
				svc:  nil,
				repo: nil,
				handlers: []hcli.Handler{
					func() hcli.Handler {
						m := _mock.NewHandler(t)
						m.On("Name").Return("noop").Maybe()
						m.On("Exec").Return(nil).Maybe()
						return m
					}(),
				},
			},
			args: args{
				eg: func() *errgroup.Group {
					eg, _ := errgroup.WithContext(context.Background())
					return eg
				}(),
			},
		},
		{
			name: "Error",
			fields: fields{
				ctx: context.Background(),
				comp: &component.Components{
					Config: &cconfig.Config{
						Data: &config.Data{
							Reconciliation: reconciliation.Reconciliation{
								Action: "noop",
							},
						},
					},
				},
				svc:  nil,
				repo: nil,
				handlers: []hcli.Handler{
					func() hcli.Handler {
						m := _mock.NewHandler(t)
						m.On("Name").Return("noop").Maybe()
						m.On("Exec").Return(errors.New("error")).Maybe()
						return m
					}(),
				},
			},
			args: args{
				eg: func() *errgroup.Group {
					eg, _ := errgroup.WithContext(context.Background())
					return eg
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{
				ctx:      tt.fields.ctx,
				comp:     tt.fields.comp,
				svc:      tt.fields.svc,
				repo:     tt.fields.repo,
				handlers: tt.fields.handlers,
			}
			c.Start(tt.args.eg)
			bf.Reset()
		})
	}
}

func TestNewCli(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		comp     *component.Components
		svc      *service.Services
		repo     *repository.Repositories
		handlers []hcli.Handler
	}

	handlerMoc := func() hcli.Handler {
		m := _mock.NewHandler(t)
		m.On("SetComponents", mock.Anything).Maybe()
		m.On("SetServices", mock.Anything).Maybe()
		m.On("SetRepositories", mock.Anything).Maybe()
		return m
	}()

	logger := clogger.NewLogger(
		context.Background(),
		&bf,
	)

	tests := []struct {
		want    *Cli
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				comp: &component.Components{
					Logger: logger,
				},
				svc:  nil,
				repo: nil,
				handlers: []hcli.Handler{
					handlerMoc,
				},
			},
			want: &Cli{
				ctx: logger.GetCtx(),
				comp: &component.Components{
					Logger: logger,
				},
				svc:  nil,
				repo: nil,
				handlers: []hcli.Handler{
					handlerMoc,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCli(tt.args.comp, tt.args.svc, tt.args.repo, tt.args.handlers)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCli() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCli() got = %v, want %v", got, tt.want)
			}
		})
	}
}

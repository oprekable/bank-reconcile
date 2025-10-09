package appcontext

import (
	"bytes"
	"context"
	"embed"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"
	"github.com/stretchr/testify/assert"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/server/cli"
	"github.com/oprekable/bank-reconcile/internal/app/service"

	"golang.org/x/sync/errgroup"
)

// checkPprofFile checks if input pprof files exist
func checkPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		info, err := os.Stat(pprof)
		assert.Nil(t, err)
		assert.False(t, os.IsNotExist(err))
		assert.False(t, info.IsDir())
	}
}

// cleanupPprofFiles deletes all specified pprof files
func cleanupPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		err := os.Remove(pprof)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAppContextGetComponents(t *testing.T) {
	type fields struct {
		ctx          context.Context
		ctxCancel    context.CancelFunc
		eg           *errgroup.Group
		embedFS      *embed.FS
		repositories *repository.Repositories
		services     *service.Services
		components   *component.Components
		servers      *server.Server
	}

	tests := []struct {
		fields fields
		want   *component.Components
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:          nil,
				ctxCancel:    nil,
				eg:           nil,
				embedFS:      nil,
				repositories: nil,
				services:     nil,
				components:   &component.Components{},
				servers:      nil,
			},
			want: &component.Components{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppContext{
				ctx:          tt.fields.ctx,
				ctxCancel:    tt.fields.ctxCancel,
				eg:           tt.fields.eg,
				embedFS:      tt.fields.embedFS,
				repositories: tt.fields.repositories,
				services:     tt.fields.services,
				components:   tt.fields.components,
				servers:      tt.fields.servers,
			}

			if got := a.GetComponents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppContextGetCtx(t *testing.T) {
	type fields struct {
		ctx          context.Context
		ctxCancel    context.CancelFunc
		eg           *errgroup.Group
		embedFS      *embed.FS
		repositories *repository.Repositories
		services     *service.Services
		components   *component.Components
		servers      *server.Server
	}

	var bf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	logger := clogger.NewLogger(
		ctx,
		&bf,
	)

	tests := []struct {
		fields fields
		want   context.Context
		name   string
	}{
		{
			name: "Ctx from component",
			fields: fields{
				ctx:          ctx,
				ctxCancel:    cancel,
				eg:           eg,
				embedFS:      nil,
				repositories: nil,
				services:     nil,
				components: &component.Components{
					Logger: logger,
				},
				servers: nil,
			},
			want: logger.GetCtx(),
		},
		{
			name: "Ctx from normal ctx",
			fields: fields{
				ctx:          ctx,
				ctxCancel:    cancel,
				eg:           eg,
				embedFS:      nil,
				repositories: nil,
				services:     nil,
				components:   nil,
				servers:      nil,
			},
			want: ctx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppContext{
				ctx:          tt.fields.ctx,
				ctxCancel:    tt.fields.ctxCancel,
				eg:           tt.fields.eg,
				embedFS:      tt.fields.embedFS,
				repositories: tt.fields.repositories,
				services:     tt.fields.services,
				components:   tt.fields.components,
				servers:      tt.fields.servers,
			}

			if got := a.GetCtx(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppContextStart(t *testing.T) {
	type fields struct {
		ctx          context.Context
		ctxCancel    context.CancelFunc
		eg           *errgroup.Group
		embedFS      *embed.FS
		repositories *repository.Repositories
		services     *service.Services
		components   *component.Components
		servers      *server.Server
	}

	var bf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	logger := clogger.NewLogger(
		ctx,
		&bf,
	)

	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				ctx:          ctx,
				ctxCancel:    cancel,
				eg:           eg,
				embedFS:      nil,
				repositories: nil,
				services:     nil,
				components: &component.Components{
					Logger: logger,
					Config: &cconfig.Config{
						Data: &config.Data{
							App: core.App{
								IsProfilerActive: true,
							},
						},
					},
					Profiler: cprofiler.NewProfiler(logger),
				},
				servers: server.NewServer(
					func() server.IServer {
						m, _ := cli.NewCli(
							&component.Components{
								Logger: logger,
								Config: &cconfig.Config{
									Data: &config.Data{
										Reconciliation: reconciliation.Reconciliation{
											Action: "noop",
										},
									},
								},
							},
							nil,
							nil,
							[]hcli.Handler{
								noop.NewHandler(&bf),
							},
						)
						return m
					}(),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppContext{
				ctx:          tt.fields.ctx,
				ctxCancel:    tt.fields.ctxCancel,
				eg:           tt.fields.eg,
				embedFS:      tt.fields.embedFS,
				repositories: tt.fields.repositories,
				services:     tt.fields.services,
				components:   tt.fields.components,
				servers:      tt.fields.servers,
			}
			go func() {
				time.Sleep(100 * time.Millisecond)
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			}()

			er := a.Start()
			t.Log(er)
			bf.Reset()

			checkPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})

			cleanupPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})
		})
	}
}

func TestNewAppContext(t *testing.T) {
	type args struct {
		ctx        context.Context
		embedFS    *embed.FS
		repository *repository.Repositories
		services   *service.Services
		components *component.Components
		servers    *server.Server
	}

	tests := []struct {
		args args
		want *AppContext
		name string
	}{
		{
			name: "Ok",
			args: args{
				ctx:        context.Background(),
				embedFS:    &embed.FS{},
				repository: &repository.Repositories{},
				services:   &service.Services{},
				components: &component.Components{},
				servers:    &server.Server{},
			},
			want: &AppContext{
				embedFS:      &embed.FS{},
				repositories: &repository.Repositories{},
				services:     &service.Services{},
				components:   &component.Components{},
				servers:      &server.Server{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewAppContext(tt.args.ctx, tt.args.embedFS, tt.args.repository, tt.args.services, tt.args.components, tt.args.servers)

			if !reflect.DeepEqual(got.embedFS, tt.want.embedFS) {
				t.Errorf("NewAppContext() got.embedFS = %v, want.embedFS %v", got.embedFS, tt.want.embedFS)
			}

			if !reflect.DeepEqual(got.repositories, tt.want.repositories) {
				t.Errorf("NewAppContext() got.repositories = %v, want.repositories %v", got.repositories, tt.want.repositories)
			}

			if !reflect.DeepEqual(got.services, tt.want.services) {
				t.Errorf("NewAppContext() got.services = %v, want.services %v", got.services, tt.want.services)
			}

			if !reflect.DeepEqual(got.components, tt.want.components) {
				t.Errorf("NewAppContext() got.components = %v, want.components %v", got.components, tt.want.components)
			}

			if !reflect.DeepEqual(got.servers, tt.want.servers) {
				t.Errorf("NewAppContext() got.servers = %v, want.servers %v", got.servers, tt.want.servers)
			}
		})
	}
}

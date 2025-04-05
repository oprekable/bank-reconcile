package helper

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	ccore "github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/app/err/core"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/server/cli"
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func TestNewRunner(t *testing.T) {
	type args struct {
		wireApp WireApp
		cmd     *cobra.Command
		args    []string
	}

	tests := []struct {
		name string
		args args
		want *Runner
	}{
		{
			name: "Ok",
			args: args{
				wireApp: nil,
				cmd:     nil,
				args:    nil,
			},
			want: &Runner{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRunner(tt.args.wireApp, tt.args.cmd, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRunner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunnerRun(t *testing.T) {
	var bf bytes.Buffer
	ctx := context.Background()
	logger := clogger.NewLogger(
		ctx,
		&bf,
	)

	type fields struct {
		wireApp WireApp
		cmd     *cobra.Command
		args    []string
	}

	type args struct {
		embedFs   *embed.FS
		appName   string
		tz        string
		errTypes  []core.ErrorType
		isVerbose bool
		dbPath    csqlite.DBPath
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				wireApp: func(
					ctx context.Context,
					embedFS *embed.FS,
					appName cconfig.AppName,
					tz cconfig.TimeZone,
					errType []core.ErrorType,
					isShowLog clogger.IsShowLog,
					dBPath csqlite.DBPath,
				) (returnData *appcontext.AppContext, cancel func(), err error) {
					returnData, cancel = appcontext.NewAppContext(
						ctx,
						nil,
						nil,
						nil,
						&component.Components{
							Config: &cconfig.Config{
								Data: &config.Data{
									App:            ccore.App{},
									Reconciliation: reconciliation.Reconciliation{},
								},
							},
						},
						server.NewServer(
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
					)

					root.FlagToDateValue = "2025-03-05"
					root.FlagFromDateValue = "2025-03-05"

					return
				},
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				args: nil,
			},
			args:    args{},
			wantErr: false,
		},
		{
			name: "Error wireApp",
			fields: fields{
				wireApp: func(
					ctx context.Context,
					embedFS *embed.FS,
					appName cconfig.AppName,
					tz cconfig.TimeZone,
					errType []core.ErrorType,
					isShowLog clogger.IsShowLog,
					dBPath csqlite.DBPath,
				) (returnData *appcontext.AppContext, cancel func(), err error) {
					err = errors.New("error")
					return
				},
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				args: nil,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "Err parse FlagToDateValue",
			fields: fields{
				wireApp: func(
					ctx context.Context,
					embedFS *embed.FS,
					appName cconfig.AppName,
					tz cconfig.TimeZone,
					errType []core.ErrorType,
					isShowLog clogger.IsShowLog,
					dBPath csqlite.DBPath,
				) (returnData *appcontext.AppContext, cancel func(), err error) {
					returnData, cancel = appcontext.NewAppContext(
						ctx,
						nil,
						nil,
						nil,
						&component.Components{
							Config: &cconfig.Config{
								Data: &config.Data{
									App:            ccore.App{},
									Reconciliation: reconciliation.Reconciliation{},
								},
							},
						},
						server.NewServer(
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
					)

					root.FlagToDateValue = "random"

					return
				},
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				args: nil,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "Err parse FlagFromDateValue",
			fields: fields{
				wireApp: func(
					ctx context.Context,
					embedFS *embed.FS,
					appName cconfig.AppName,
					tz cconfig.TimeZone,
					errType []core.ErrorType,
					isShowLog clogger.IsShowLog,
					dBPath csqlite.DBPath,
				) (returnData *appcontext.AppContext, cancel func(), err error) {
					returnData, cancel = appcontext.NewAppContext(
						ctx,
						nil,
						nil,
						nil,
						&component.Components{
							Config: &cconfig.Config{
								Data: &config.Data{
									App:            ccore.App{},
									Reconciliation: reconciliation.Reconciliation{},
								},
							},
						},
						server.NewServer(
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
					)

					root.FlagToDateValue = "2025-03-05"
					root.FlagFromDateValue = "random"

					return
				},
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				args: nil,
			},
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				wireApp: tt.fields.wireApp,
				cmd:     tt.fields.cmd,
				args:    tt.fields.args,
			}

			if err := r.Run(tt.args.embedFs, tt.args.appName, tt.args.tz, tt.args.errTypes, tt.args.isVerbose, tt.args.dbPath); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			bf.Reset()
		})
	}
}

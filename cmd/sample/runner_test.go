package sample

import (
	"bytes"
	"context"
	"embed"
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
	"testing"
)

func TestRunner(t *testing.T) {
	var bf bytes.Buffer
	ctx := context.Background()
	logger := clogger.NewLogger(
		ctx,
		&bf,
	)

	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name        string
		args        args
		triggerMock func()
		wantErr     bool
	}{
		{
			name: "Ok",
			args: args{
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				args: []string{},
			},
			triggerMock: func() {
				wireApp = func(
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
				}

				root.FlagIsDebugValue = true
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.triggerMock()

			if err := Runner(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}

			bf.Reset()
		})
	}
}

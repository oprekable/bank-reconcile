package process

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	core2 "github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/app/err/core"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/noop"
	"github.com/oprekable/bank-reconcile/internal/app/server"
	"github.com/oprekable/bank-reconcile/internal/app/server/cli"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/spf13/cobra"
)

var wireApp = func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
	return &appcontext.AppContext{}, nil, nil
}

func TestCmdProcessInit(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		appName      string
	}

	type args struct {
		in0 *cmd.MetaData
	}

	tests := []struct {
		fields fields
		args   args
		want   *cobra.Command
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				c: &cobra.Command{
					Use:     Usage,
					Short:   Short,
					Long:    Long,
					Aliases: Aliases,
					Example: fmt.Sprintf(
						"%s\n",
						fmt.Sprintf("Process data \n\t%s %s", "", Example),
					),
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			args: args{
				in0: nil,
			},
			want: func() *cobra.Command {
				c := NewCommand("", wireApp, nil, &bytes.Buffer{}, &bytes.Buffer{})
				c.c.Use = Usage
				c.c.Short = Short
				c.c.Long = Long
				c.c.RunE = c.Runner

				c.c.Example = fmt.Sprintf(
					"%s\n",
					fmt.Sprintf("Process data \n\t%s %s", c.appName, Example),
				)

				c.c.SetOut(&bytes.Buffer{})
				c.c.SetErr(&bytes.Buffer{})

				return c.c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdProcess{
				c:            tt.fields.c,
				appName:      tt.fields.appName,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			got := c.Init(tt.args.in0)

			if !reflect.DeepEqual(got.Use, tt.want.Use) ||
				!reflect.DeepEqual(got.Short, tt.want.Short) ||
				!reflect.DeepEqual(got.Long, tt.want.Long) ||
				!reflect.DeepEqual(got.Example, tt.want.Example) ||
				!reflect.DeepEqual(got.OutOrStdout(), tt.want.OutOrStdout()) ||
				!reflect.DeepEqual(got.OutOrStderr(), tt.want.OutOrStderr()) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}

			got.OutOrStdout().(*bytes.Buffer).Reset()
			got.OutOrStderr().(*bytes.Buffer).Reset()
		})
	}
}

func TestCmdProcessPersistentPreRunner(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		appName      string
	}

	type args struct {
		cCmd *cobra.Command
		args []string
	}

	tests := []struct {
		fields  fields
		trigger func()
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				cCmd: nil,
				args: nil,
			},
			trigger: func() {
				cmd.FlagTotalDataSampleToGenerateValue = 1000
				cmd.FlagPercentageMatchSampleToGenerateValue = 100
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: false,
		},
		{
			name: "Error - invalid from date",
			args: args{
				cCmd: nil,
				args: nil,
			},
			trigger: func() {
				cmd.FlagTotalDataSampleToGenerateValue = 1000
				cmd.FlagPercentageMatchSampleToGenerateValue = 100
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-xx"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdProcess{
				c:            tt.fields.c,
				appName:      tt.fields.appName,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			tt.trigger()

			if err := c.PersistentPreRunner(tt.args.cCmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCmdProcessRunner(t *testing.T) {
	var bf bytes.Buffer
	ctx := context.Background()
	logger := clogger.NewLogger(
		ctx,
		&bf,
	)

	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		appName      string
	}

	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		fields  fields
		trigger func()
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				c: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				appName: "",
				wireApp: func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
					app, cancel := appcontext.NewAppContext(
						ctx,
						nil,
						nil,
						nil,
						&component.Components{
							Logger: logger,
							Config: &cconfig.Config{
								Data: &config.Data{
									App:            core2.App{},
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

					return app, cancel, nil
				},
				embedFS:      nil,
				outPutWriter: nil,
				errWriter:    nil,
			},
			args: args{},
			trigger: func() {
				cmd.FlagIsVerboseValue = true
				cmd.FlagIsDebugValue = true
				cmd.FlagIsProfilerActiveValue = true
				cmd.FlagSystemTRXPathValue = "/tmp/sample/system"
				cmd.FlagBankTRXPathValue = "/tmp/sample/bank"
				cmd.FlagReportTRXPathValue = "/tmp/report"
				cmd.FlagListBankValue = []string{"foo", "bar"}
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: false,
		},
		{
			name: "Error - dependency injection cause error",
			fields: fields{
				c: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				appName: "",
				wireApp: func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
					return nil, nil, errors.New("dependency-injection error")
				},
				embedFS:      nil,
				outPutWriter: nil,
				errWriter:    nil,
			},
			args:    args{},
			trigger: func() {},
			wantErr: true,
		},
		{
			name: "Error - simulate invalid flag value (fromDate)",
			fields: fields{
				c: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				appName: "",
				wireApp: func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
					app, cancel := appcontext.NewAppContext(
						ctx,
						nil,
						nil,
						nil,
						&component.Components{
							Logger: logger,
							Config: &cconfig.Config{
								Data: &config.Data{
									App:            core2.App{},
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

					return app, cancel, nil
				},
				embedFS:      nil,
				outPutWriter: nil,
				errWriter:    nil,
			},
			args: args{},
			trigger: func() {
				cmd.FlagIsVerboseValue = true
				cmd.FlagIsDebugValue = true
				cmd.FlagIsProfilerActiveValue = true
				cmd.FlagSystemTRXPathValue = "/tmp/sample/system"
				cmd.FlagBankTRXPathValue = "/tmp/sample/bank"
				cmd.FlagReportTRXPathValue = "/tmp/report"
				cmd.FlagListBankValue = []string{"foo", "bar"}
				cmd.FlagFromDateValue = "2025-05-xx"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdProcess{
				c:            tt.fields.c,
				appName:      tt.fields.appName,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			tt.trigger()

			if err := c.Runner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}

			bf.Reset()
		})
	}
}

func TestCmdProcessInitPersistentFlags(t *testing.T) {
	wD := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
	dateNow := time.Now().Format("2006-01-02")
	bf := &bytes.Buffer{}

	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		appName      string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				c: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetOut(bf)
					r.SetErr(bf)
					return r
				}(),
				outPutWriter: bf,
				errWriter:    bf,
			},
			want: fmt.Sprintf(`  -b, --banktrxpath string     Path location of Bank Transaction directory (default "%s/sample/bank")
  -g, --debug                  debug mode
  -d, --deleteoldfile          delete old report files (default true)
  -f, --from string            from date (YYYY-MM-DD) (default "%s")
  -l, --listbank strings       List bank accepted (default [bca,bni,mandiri,bri,danamon])
  -i, --profiler               pprof active mode
  -r, --reportpath string      Path location of Archive directory (default "%s/report")
  -o, --showlog                show logs
  -s, --systemtrxpath string   Path location of System Transaction directory (default "%s/sample/system")
  -z, --time_zone string       time zone settings (default "Asia/Jakarta")
  -t, --to string              to date (YYYY-MM-DD) (default "%s")
`,
				wD,
				dateNow,
				wD,
				wD,
				dateNow,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdProcess{
				c:            tt.fields.c,
				appName:      tt.fields.appName,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			c.initPersistentFlags()
			got := tt.fields.c.PersistentFlags().FlagUsages()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitCommonPersistentFlags() = %v, want %v", got, tt.want)
			}

			c.outPutWriter.(*bytes.Buffer).Reset()
			c.errWriter.(*bytes.Buffer).Reset()
		})
	}
}

func TestNewCommand(t *testing.T) {
	type args struct {
		wireApp inject.Fn
		embedFS *embed.FS
		appName string
	}

	tests := []struct {
		args             args
		want             *CmdProcess
		name             string
		wantOutPutWriter string
		wantErrWriter    string
	}{
		{
			name: "Ok",
			args: args{
				wireApp: wireApp,
				embedFS: nil,
			},
			wantOutPutWriter: "",
			wantErrWriter:    "",
			want: &CmdProcess{
				c: &cobra.Command{
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPutWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}
			got := NewCommand(tt.args.appName, tt.args.wireApp, tt.args.embedFS, outPutWriter, errWriter)
			if gotOutPutWriter := outPutWriter.String(); gotOutPutWriter != tt.wantOutPutWriter {
				t.Errorf("NewCommand() gotOutPutWriter = %v, want %v", gotOutPutWriter, tt.wantOutPutWriter)
			}

			if gotErrWriter := errWriter.String(); gotErrWriter != tt.wantErrWriter {
				t.Errorf("NewCommand() gotErrWriter = %v, want %v", gotErrWriter, tt.wantErrWriter)
			}

			if !reflect.DeepEqual(got.outPutWriter, tt.want.outPutWriter) ||
				!reflect.DeepEqual(got.errWriter, tt.want.errWriter) ||
				!reflect.DeepEqual(got.c.SilenceErrors, tt.want.c.SilenceErrors) ||
				!reflect.DeepEqual(got.c.SilenceUsage, tt.want.c.SilenceUsage) ||
				!reflect.DeepEqual(reflect.ValueOf(got.wireApp).Pointer(), reflect.ValueOf(tt.want.wireApp).Pointer()) {
				t.Errorf("NewCommand() = %v, want %v", got, tt.want)
			}

			outPutWriter.Reset()
			errWriter.Reset()
		})
	}
}

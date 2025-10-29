package sample

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/cmd/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"
	"github.com/stretchr/testify/mock"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/_inject"
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
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/spf13/cobra"
)

var wireApp = func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
	return &appcontext.AppContext{}, nil, nil
}

// cleanupPprofFiles deletes all specified pprof files
func cleanupPprofFiles(_ *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		_ = os.Remove(pprof)
	}
}

func TestCmdSampleInit(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      _inject.Fn
		embedFS      *embed.FS
		appName      string
		subCommands  []cmd.Cmd
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
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
				subCommands: []cmd.Cmd{
					func() cmd.Cmd {
						m := _mock.NewCmd(t)
						m.On(
							"Init",
							mock.Anything,
						).Return(&cobra.Command{}).
							Maybe()

						m.On(
							"Example",
						).Return("example string").
							Maybe()

						return m
					}(),
				},
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
					fmt.Sprintf("Generate sample \n\t%s %s", c.appName, Example),
				)

				c.c.SetOut(&bytes.Buffer{})
				c.c.SetErr(&bytes.Buffer{})

				return c.c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdSample{
				c:            tt.fields.c,
				appName:      tt.fields.appName,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
				subCommands:  tt.fields.subCommands,
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

func TestCmdSamplePersistentPreRunner(t *testing.T) {
	type args struct {
		cCmd *cobra.Command
		args []string
	}

	tests := []struct {
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
			name: "Error - invalid total data",
			args: args{
				cCmd: nil,
				args: nil,
			},
			trigger: func() {
				cmd.FlagTotalDataSampleToGenerateValue = 0
				cmd.FlagPercentageMatchSampleToGenerateValue = 100
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: true,
		},
		{
			name: "Error - invalid percentage match",
			args: args{
				cCmd: nil,
				args: nil,
			},
			trigger: func() {
				cmd.FlagTotalDataSampleToGenerateValue = 1000
				cmd.FlagPercentageMatchSampleToGenerateValue = -10
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdSample{}
			tt.trigger()

			if err := c.PersistentPreRunner(tt.args.cCmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCmdSampleRunner(t *testing.T) {
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
		wireApp      _inject.Fn
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
							Profiler: cprofiler.NewProfiler(logger),
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
				cmd.FlagListBankValue = []string{"foo", "bar"}
				cmd.FlagFromDateValue = "2025-05-xx"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdSample{
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

			cleanupPprofFiles(t, []string{
				"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
				"./trace.pprof", "./goroutine.pprof",
			})
		})
	}
}

func TestCmdSampleInitPersistentFlags(t *testing.T) {
	wD := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
	dateNow := time.Now().Format("2006-01-02")
	bf := &bytes.Buffer{}

	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      _inject.Fn
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
			want: fmt.Sprintf(`  -a, --amountdata int         amount system trx data sample to generate, bank trx will be 2 times of this amount (default 1000)
  -b, --banktrxpath string     Path location of Bank Transaction directory (default "%s/sample/bank")
  -g, --debug                  debug mode
  -d, --deleteoldfile          delete old sample files (default true)
  -f, --from string            from date (YYYY-MM-DD) (default "%s")
  -l, --listbank strings       List bank accepted (default [bca,bni,mandiri,bri,danamon])
  -p, --percentagematch int    percentage of matched trx for data sample to generate (default 100)
  -i, --profiler               pprof active mode
  -o, --showlog                show logs
  -s, --systemtrxpath string   Path location of System Transaction directory (default "%s/sample/system")
  -z, --time_zone string       time zone settings (default "Asia/Jakarta")
  -t, --to string              to date (YYYY-MM-DD) (default "%s")
`,
				wD,
				dateNow,
				wD,
				dateNow,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdSample{
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
		wireApp _inject.Fn
		embedFS *embed.FS
		appName string
	}

	tests := []struct {
		args             args
		want             *CmdSample
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
			want: &CmdSample{
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

func TestCmdSampleExample(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		wireApp      _inject.Fn
		embedFS      *embed.FS
		appName      string
		subCommands  []cmd.Cmd
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				c: &cobra.Command{
					Use:           Usage,
					Short:         Short,
					Long:          Long,
					Aliases:       Aliases,
					Example:       "example string",
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			want: "example string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdSample{
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
				c:            tt.fields.c,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				appName:      tt.fields.appName,
				subCommands:  tt.fields.subCommands,
			}
			if got := c.Example(); got != tt.want {
				t.Errorf("Example() = %v, want %v", got, tt.want)
			}
		})
	}
}

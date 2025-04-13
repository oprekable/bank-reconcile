package helper

import (
	"bytes"
	"fmt"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/spf13/cobra"
	"reflect"
	"testing"
	"time"
)

func TestCommonPersistentPreRunner(t *testing.T) {
	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		name    string
		args    args
		trigger func()
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				in0: nil,
				in1: nil,
			},
			trigger: func() {
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-01"
			},
			wantErr: false,
		},
		{
			name: "Error - toDate less than fromDate",
			args: args{
				in0: nil,
				in1: nil,
			},
			trigger: func() {
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-04-01"
			},
			wantErr: true,
		},
		{
			name: "Error - invalid toDate",
			args: args{
				in0: nil,
				in1: nil,
			},
			trigger: func() {
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-04-xx"
			},
			wantErr: true,
		},
		{
			name: "Error - invalid fromDate",
			args: args{
				in0: nil,
				in1: nil,
			},
			trigger: func() {
				cmd.FlagTZValue = time.UTC.String()
				cmd.FlagFromDateValue = "2025-05-xx"
			},
			wantErr: true,
		},
		{
			name: "Error - invalid timezone",
			args: args{
				in0: nil,
				in1: nil,
			},
			trigger: func() {
				cmd.FlagTZValue = "any string"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.trigger()

			if err := CommonPersistentPreRunner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("CommonPersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitCommonPersistentFlags(t *testing.T) {
	workDir := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
	dateNow := time.Now().Format("2006-01-02")
	bf := &bytes.Buffer{}
	type args struct {
		c *cobra.Command
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Ok",
			args: args{
				c: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetOut(bf)
					r.SetErr(bf)
					return r
				}(),
			},
			want: fmt.Sprintf(`  -b, --banktrxpath string     Path location of Bank Transaction directory (default "%s/sample/bank")
  -g, --debug                  debug mode
  -f, --from string            from date (YYYY-MM-DD) (default "%s")
  -l, --listbank strings       List bank accepted (default [bca,bni,mandiri,bri,danamon])
  -i, --profiler               pprof active mode
  -o, --showlog                show logs
  -s, --systemtrxpath string   Path location of System Transaction directory (default "%s/sample/system")
  -z, --time_zone string       time zone settings (default "Asia/Jakarta")
  -t, --to string              to date (YYYY-MM-DD) (default "%s")
`,
				workDir,
				dateNow,
				workDir,
				dateNow,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCommonPersistentFlags(tt.args.c)
			got := tt.args.c.PersistentFlags().FlagUsages()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitCommonPersistentFlags() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestUpdateCommonConfigFromFlags(t *testing.T) {
	type args struct {
		app appcontext.IAppContext
	}

	tests := []struct {
		name    string
		args    args
		trigger func()
		want    *cconfig.Config
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				app: func() appcontext.IAppContext {
					r := _mock.NewIAppContext(t)
					r.On("GetComponents").
						Return(
							&component.Components{
								Config: &cconfig.Config{
									Data: &config.Data{
										Reconciliation: reconciliation.Reconciliation{},
									},
								},
							},
						).
						Maybe()
					return r
				}(),
			},
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
			want: &cconfig.Config{
				Data: &config.Data{
					App: core.App{
						IsShowLog:        true,
						IsDebug:          true,
						IsProfilerActive: true,
					},
					Reconciliation: reconciliation.Reconciliation{
						FromDate: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-05-01")
							return t
						}(),
						ToDate: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-05-01")
							return t
						}(),
						SystemTRXPath: "/tmp/sample/system",
						BankTRXPath:   "/tmp/sample/bank",
						ListBank:      []string{"foo", "bar"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error - toDate",
			args: args{
				app: func() appcontext.IAppContext {
					r := _mock.NewIAppContext(t)
					r.On("GetComponents").
						Return(
							&component.Components{
								Config: &cconfig.Config{
									Data: &config.Data{
										Reconciliation: reconciliation.Reconciliation{},
									},
								},
							},
						).
						Maybe()
					return r
				}(),
			},
			trigger: func() {
				cmd.FlagIsVerboseValue = true
				cmd.FlagIsDebugValue = true
				cmd.FlagIsProfilerActiveValue = true
				cmd.FlagSystemTRXPathValue = "/tmp/sample/system"
				cmd.FlagBankTRXPathValue = "/tmp/sample/bank"
				cmd.FlagListBankValue = []string{"foo", "bar"}
				cmd.FlagFromDateValue = "2025-05-01"
				cmd.FlagToDateValue = "2025-05-xx"
			},
			want: &cconfig.Config{
				Data: &config.Data{
					App: core.App{
						IsShowLog:        true,
						IsDebug:          true,
						IsProfilerActive: true,
					},
					Reconciliation: reconciliation.Reconciliation{
						FromDate:      time.Time{},
						ToDate:        time.Time{},
						SystemTRXPath: "/tmp/sample/system",
						BankTRXPath:   "/tmp/sample/bank",
						ListBank:      []string{"foo", "bar"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Ok",
			args: args{
				app: func() appcontext.IAppContext {
					r := _mock.NewIAppContext(t)
					r.On("GetComponents").
						Return(
							&component.Components{
								Config: &cconfig.Config{
									Data: &config.Data{
										Reconciliation: reconciliation.Reconciliation{},
									},
								},
							},
						).
						Maybe()
					return r
				}(),
			},
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
			want: &cconfig.Config{
				Data: &config.Data{
					App: core.App{
						IsShowLog:        true,
						IsDebug:          true,
						IsProfilerActive: true,
					},
					Reconciliation: reconciliation.Reconciliation{
						FromDate: time.Time{},
						ToDate: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-05-01")
							return t
						}(),
						SystemTRXPath: "/tmp/sample/system",
						BankTRXPath:   "/tmp/sample/bank",
						ListBank:      []string{"foo", "bar"},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.trigger()

			if err := UpdateCommonConfigFromFlags(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCommonConfigFromFlags() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := tt.args.app.GetComponents().Config
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateCommonConfigFromFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

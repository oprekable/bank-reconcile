package root

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/_mock"
	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/internal/app/appcontext"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/err/core"
	"github.com/oprekable/bank-reconcile/internal/inject"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"io"
	"reflect"
	"testing"
	"time"
)

var wireApp = func(ctx context.Context, embedFS *embed.FS, appName cconfig.AppName, tz cconfig.TimeZone, errType []core.ErrorType, isShowLog clogger.IsShowLog, dBPath csqlite.DBPath) (*appcontext.AppContext, func(), error) {
	return &appcontext.AppContext{}, nil, nil
}

func TestCmdRootInit(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		outPutWriter io.Writer
		errWriter    io.Writer
		subCommands  []cmd.Cmd
	}

	type args struct {
		metaData *cmd.MetaData
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *cobra.Command
	}{
		{
			name: "Ok",
			fields: fields{
				c: &cobra.Command{
					Args:          cobra.NoArgs,
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
				subCommands: func() []cmd.Cmd {
					m := _mock.NewCmd(t)
					m.On(
						"Init",
						mock.Anything,
					).Return(&cobra.Command{}).
						Maybe()

					return []cmd.Cmd{
						m,
					}
				}(),
			},
			args: args{
				metaData: &cmd.MetaData{
					Usage: "foo",
					Short: "f",
					Long:  "foo foo foo",
				},
			},
			want: func() *cobra.Command {
				c := NewCommand(wireApp, nil, &bytes.Buffer{}, &bytes.Buffer{})

				c.c.Use = "foo"
				c.c.Short = "f"
				c.c.Long = "foo foo foo"
				c.c.RunE = c.Runner

				c.c.Example = fmt.Sprintf(
					"%s\n%s\n",
					fmt.Sprintf("Generate sample \n\t%s %s", "foo", sample.Example),
					fmt.Sprintf("Process data \n\t%s %s", "foo", process.Example),
				)

				c.c.SetOut(&bytes.Buffer{})
				c.c.SetErr(&bytes.Buffer{})

				return c.c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdRoot{
				c:            tt.fields.c,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
				subCommands:  tt.fields.subCommands,
			}

			got := c.Init(tt.args.metaData)

			if !reflect.DeepEqual(got.Use, tt.want.Use) ||
				!reflect.DeepEqual(got.Short, tt.want.Short) ||
				!reflect.DeepEqual(got.Long, tt.want.Long) ||
				!reflect.DeepEqual(got.Example, tt.want.Example) ||
				!reflect.DeepEqual(got.OutOrStdout(), tt.want.OutOrStdout()) ||
				!reflect.DeepEqual(got.OutOrStderr(), tt.want.OutOrStderr()) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCmdRootPersistentPreRunner(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		outPutWriter io.Writer
		errWriter    io.Writer
		subCommands  []cmd.Cmd
	}

	type args struct {
		in0 *cobra.Command
		in1 []string
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
				c:            nil,
				wireApp:      nil,
				embedFS:      nil,
				outPutWriter: nil,
				errWriter:    nil,
				subCommands:  nil,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdRoot{
				c:            tt.fields.c,
				wireApp:      tt.fields.wireApp,
				embedFS:      tt.fields.embedFS,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
				subCommands:  tt.fields.subCommands,
			}

			if err := c.PersistentPreRunner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCmdRootRunner(t *testing.T) {
	workDir := filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
	dateNow := time.Now().Format("2006-01-02")

	type fields struct {
		c            *cobra.Command
		wireApp      inject.Fn
		embedFS      *embed.FS
		outPutWriter io.Writer
		errWriter    io.Writer
		subCommands  []cmd.Cmd
	}

	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				c: &cobra.Command{
					Args:          cobra.NoArgs,
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
				subCommands:  nil,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			want: fmt.Sprintf(`foo foo foo

Usage:
  foo

Examples:
Generate sample 
	foo sample --systemtrxpath=%s/sample/system --banktrxpath=%s/sample/bank --listbank=bca,bni,mandiri,bri,danamon --percentagematch=100 --amountdata=10000 --from=%s --to=%s
Process data 
	foo process --systemtrxpath=%s/sample/system --banktrxpath=%s/sample/bank --reportpath==%s/report --listbank=bca,bni,mandiri,bri,danamon --from=%s --to=%s

`,
				workDir,
				workDir,
				dateNow,
				dateNow,
				workDir,
				workDir,
				workDir,
				dateNow,
				dateNow,
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommand(
				tt.fields.wireApp,
				tt.fields.embedFS,
				tt.fields.outPutWriter,
				tt.fields.errWriter,
			)

			c.Init(
				&cmd.MetaData{
					Usage: "foo",
					Short: "f",
					Long:  "foo foo foo",
				},
			)

			if err := c.Runner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := c.outPutWriter.(*bytes.Buffer).String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Runner() output = %v, want %v", got, tt.want)
			}

			c.outPutWriter.(*bytes.Buffer).Reset()
			c.errWriter.(*bytes.Buffer).Reset()
		})
	}
}

func TestNewCommand(t *testing.T) {
	type args struct {
		wireApp     inject.Fn
		embedFS     *embed.FS
		subCommands []cmd.Cmd
	}

	tests := []struct {
		name             string
		args             args
		wantOutPutWriter string
		wantErrWriter    string
		want             *CmdRoot
	}{
		{
			name: "Ok",
			args: args{
				wireApp:     wireApp,
				embedFS:     nil,
				subCommands: nil,
			},
			wantOutPutWriter: "",
			wantErrWriter:    "",
			want: &CmdRoot{
				c: &cobra.Command{
					Args:          cobra.NoArgs,
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				wireApp:      wireApp,
				embedFS:      nil,
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
				subCommands:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPutWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}
			got := NewCommand(tt.args.wireApp, tt.args.embedFS, outPutWriter, errWriter, tt.args.subCommands...)

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
				!reflect.DeepEqual(reflect.ValueOf(got.c.Args).Pointer(), reflect.ValueOf(tt.want.c.Args).Pointer()) ||
				!reflect.DeepEqual(reflect.ValueOf(got.wireApp).Pointer(), reflect.ValueOf(tt.want.wireApp).Pointer()) {
				t.Errorf("NewCommand() = %v, want %v", got, tt.want)
			}

			outPutWriter.Reset()
			errWriter.Reset()
		})
	}
}

package version

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/variable"
	"github.com/spf13/cobra"
)

func TestCmdVersionInit(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		outPutWriter io.Writer
		errWriter    io.Writer
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
					Args:          cobra.NoArgs,
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			args: args{
				in0: &cmd.MetaData{},
			},
			want: func() *cobra.Command {
				c := NewCommand(&bytes.Buffer{}, &bytes.Buffer{})

				c.c.Use = Usage
				c.c.Aliases = Aliases
				c.c.Short = Short
				c.c.Long = Long
				c.c.RunE = c.Runner

				c.c.SetOut(&bytes.Buffer{})
				c.c.SetErr(&bytes.Buffer{})

				return c.c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdVersion{
				c:            tt.fields.c,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			got := c.Init(tt.args.in0)

			if !reflect.DeepEqual(got.Use, tt.want.Use) ||
				!reflect.DeepEqual(got.Aliases, tt.want.Aliases) ||
				!reflect.DeepEqual(got.Short, tt.want.Short) ||
				!reflect.DeepEqual(got.Long, tt.want.Long) ||
				!reflect.DeepEqual(got.OutOrStdout(), tt.want.OutOrStdout()) ||
				!reflect.DeepEqual(got.OutOrStderr(), tt.want.OutOrStderr()) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCmdVersionPersistentPreRunner(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		outPutWriter io.Writer
		errWriter    io.Writer
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
				outPutWriter: nil,
				errWriter:    nil,
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
			c := &CmdVersion{
				c:            tt.fields.c,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			if err := c.PersistentPreRunner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCmdVersionRunner(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		outPutWriter io.Writer
		errWriter    io.Writer
	}

	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		fields  fields
		name    string
		want    string
		args    args
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
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			want: fmt.Sprintf(
				"%s\n%s\n%s\n%s\n",
				fmt.Sprintf("App\t\t: %s", variable.AppName),
				fmt.Sprintf("Desc\t\t: %s", variable.AppDescLong),
				fmt.Sprintf("Build Date\t: %s", variable.BuildDate),
				fmt.Sprintf("Git Commit\t: %s", variable.GitCommit),
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdVersion{
				c:            tt.fields.c,
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
			}

			err := c.Runner(tt.args.in0, tt.args.in1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := c.outPutWriter.(*bytes.Buffer).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Runner() output = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		want             *CmdVersion
		name             string
		wantOutPutWriter string
		wantErrWriter    string
	}{
		{
			name:             "Ok",
			wantOutPutWriter: "",
			wantErrWriter:    "",
			want: &CmdVersion{
				c: &cobra.Command{
					Args:          cobra.NoArgs,
					SilenceErrors: true,
					SilenceUsage:  true,
				},
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPutWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			got := NewCommand(outPutWriter, errWriter)
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
				!reflect.DeepEqual(reflect.ValueOf(got.c.Args).Pointer(), reflect.ValueOf(tt.want.c.Args).Pointer()) {
				t.Errorf("NewCommand() = %v, want %v", got, tt.want)
			}

			outPutWriter.Reset()
			errWriter.Reset()
		})
	}
}

func TestCmdVersionExample(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
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
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			want: "example string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdVersion{
				outPutWriter: tt.fields.outPutWriter,
				errWriter:    tt.fields.errWriter,
				c:            tt.fields.c,
			}

			if got := c.Example(); got != tt.want {
				t.Errorf("Example() = %v, want %v", got, tt.want)
			}
		})
	}
}

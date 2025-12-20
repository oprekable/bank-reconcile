package root

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/cmd/_mock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

func TestCmdRootInit(t *testing.T) {
	type fields struct {
		c            *cobra.Command
		outPutWriter io.Writer
		errWriter    io.Writer
		subCommands  []cmd.Cmd
	}

	type args struct {
		metaData *cmd.MetaData
	}

	tests := []struct {
		args   args
		want   *cobra.Command
		name   string
		fields fields
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
				subCommands: func() []cmd.Cmd {
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
				c := NewCommand(&bytes.Buffer{}, &bytes.Buffer{})

				c.c.Use = "foo"
				c.c.Short = "f"
				c.c.Long = "foo foo foo"
				c.c.RunE = c.Runner

				c.c.Example = "example string"

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
	type fields struct {
		c            *cobra.Command
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
		want    string
		fields  fields
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
				in1: nil,
			},
			want: `foo foo foo

Usage:
  foo

Examples:
example string

Additional help topics:
  foo         
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommand(
				tt.fields.outPutWriter,
				tt.fields.errWriter,
				tt.fields.subCommands...,
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
		subCommands []cmd.Cmd
	}

	tests := []struct {
		want             *CmdRoot
		name             string
		wantOutPutWriter string
		wantErrWriter    string
		args             args
	}{
		{
			name: "Ok",
			args: args{
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
			got := NewCommand(outPutWriter, errWriter, tt.args.subCommands...)

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

func TestCmdRootExample(t *testing.T) {
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
					Example: "example string",
				},
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
			},
			want: "example string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CmdRoot{
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

func TestErrFlagsFunc(t *testing.T) {
	type fields struct {
		outPutWriter io.Writer
		errWriter    io.Writer
		c            *cobra.Command
		err          error
	}

	tests := []struct {
		fields    fields
		wantError error
		name      string
	}{
		{
			name: "Ok",
			fields: fields{
				c:            &cobra.Command{},
				outPutWriter: &bytes.Buffer{},
				errWriter:    &bytes.Buffer{},
				err:          errors.New("foo"),
			},
			wantError: errors.New("foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.c.SetErr(tt.fields.errWriter)
			tt.fields.c.SetOut(tt.fields.outPutWriter)

			gotErr := errFlagsFunc(tt.fields.c, tt.fields.err)

			if gotErr.Error() != tt.wantError.Error() {
				t.Errorf("errFlagsFunc() = err = %v, wantError %v", gotErr, tt.wantError)
			}
		})
	}
}

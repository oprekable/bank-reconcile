package helper

import (
	"bytes"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/schollz/progressbar/v3"
	"io"
	"reflect"
	"testing"
)

func TestInitCommonArgs(t *testing.T) {
	type args struct {
		conf      *config.Data
		extraArgs [][]string
	}

	tests := []struct {
		name string
		args args
		want [][]string
	}{

		{
			name: "Ok",
			args: args{
				extraArgs: [][]string{
					{
						"foo",
						"bar",
					},
				},
				conf: &config.Data{},
			},
			want: [][]string{
				{
					"-z --time_zone",
					"",
				},
				{
					"-f --from",
					"0001-01-01",
				},
				{
					"-t --to",
					"0001-01-01",
				},
				{
					"-s --systemtrxpath",
					"",
				},
				{
					"-b --banktrxpath",
					"",
				},
				{
					"-l --listbank",
					"",
				},
				{
					"-o --showlog",
					"false",
				},
				{
					"-g --debug",
					"false",
				},
				{
					"-i --profiler",
					"false",
				},
				{
					"foo",
					"bar",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitCommonArgs(tt.args.conf, tt.args.extraArgs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitCommonArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitProgressBar(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		writer io.Writer
	}

	tests := []struct {
		args args
		want *progressbar.ProgressBar
		name string
	}{
		{
			name: "Ok",
			args: args{
				writer: &bf,
			},
			want: progressbar.NewOptions(-1,
				progressbar.OptionSetWriter(&bf),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetWidth(15),
				progressbar.OptionSpinnerType(17),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitProgressBar(&bf)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitProgressBar() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

package progressbarhelper

import (
	"bytes"
	"testing"

	"github.com/schollz/progressbar/v3"
)

func TestBarClear(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		bar         *progressbar.ProgressBar
		description string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok - nil bar",
			args: args{
				bar:         nil,
				description: "foo",
			},
			want: "",
		},
		{
			name: "ok",
			args: args{
				bar:         progressbar.NewOptions(100, progressbar.OptionSetWidth(10), progressbar.OptionSetWriter(&bf)),
				description: "foo",
			},
			want: "" +
				"\rfoo   0% |          |  [0s:0s]" +
				"\r                              \r",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BarDescribe(tt.args.bar, tt.args.description)
			BarClear(tt.args.bar)
			got := bf.String()
			if got != tt.want {
				t.Errorf("BarClear() output = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

func TestBarDescribe(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		bar         *progressbar.ProgressBar
		description string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok - nil bar",
			args: args{
				bar:         nil,
				description: "foo",
			},
			want: "",
		},
		{
			name: "ok",
			args: args{
				bar:         progressbar.NewOptions(100, progressbar.OptionSetWidth(10), progressbar.OptionSetWriter(&bf)),
				description: "foo",
			},
			want: "\rfoo   0% |          |  [0s:0s]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BarDescribe(tt.args.bar, tt.args.description)
			if got := bf.String(); got != tt.want {
				t.Errorf("BarDescribe() output = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}

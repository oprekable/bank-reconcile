package version

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunner(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		triggerMock func()
		name        string
		want        string
		args        args
		wantErr     bool
	}{
		{
			name: "Ok",
			args: args{
				in0: nil,
				in1: []string{},
			},
			triggerMock: func() {
				versionWriter = &bf
			},
			want:    `App\t\t: bank-reconcile`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.triggerMock()
			if err := Runner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := bf.String()
			gotQuote := strconv.Quote(strings.TrimRight(got, "\n"))
			if !strings.Contains(gotQuote, tt.want) {
				t.Errorf("Msg() output = %v, want %v", gotQuote, tt.want)
			}

			bf.Reset()
		})
	}
}

package version

import (
	"bytes"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"testing"
)

func TestRunner(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		name        string
		args        args
		triggerMock func()
		want        string
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
			want:    `App\t\t\t: bank-reconcile`,
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

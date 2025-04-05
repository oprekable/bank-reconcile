package root

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestPersistentPreRunner(t *testing.T) {
	type args struct {
		in0 *cobra.Command
		in1 []string
	}

	tests := []struct {
		triggerMock func()
		name        string
		args        args
		wantErr     bool
	}{
		{
			name: "Ok",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "2025-03-05"
				FlagToDateValue = "2025-03-05"
				FlagTotalDataSampleToGenerateValue = 1000
				FlagPercentageMatchSampleToGenerateValue = 100
			},
			wantErr: false,
		},
		{
			name: "Error - FlagPercentageMatchSampleToGenerateValue",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "2025-03-05"
				FlagToDateValue = "2025-03-05"
				FlagTotalDataSampleToGenerateValue = 1000
				FlagPercentageMatchSampleToGenerateValue = 200
			},
			wantErr: true,
		},
		{
			name: "Error - FlagTotalDataSampleToGenerateValue",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "2025-03-05"
				FlagToDateValue = "2025-03-05"
				FlagTotalDataSampleToGenerateValue = 0
				FlagPercentageMatchSampleToGenerateValue = 100
			},
			wantErr: true,
		},
		{
			name: "Error - fromDate > toDate",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "2025-04-05"
				FlagToDateValue = "2025-03-05"
				FlagTotalDataSampleToGenerateValue = 1000
				FlagPercentageMatchSampleToGenerateValue = 100
			},
			wantErr: true,
		},
		{
			name: "Error - FlagToDateValue invalid date",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "2025-03-05"
				FlagToDateValue = "any string"
				FlagTotalDataSampleToGenerateValue = 1000
				FlagPercentageMatchSampleToGenerateValue = 100
			},
			wantErr: true,
		},
		{
			name: "Error - FlagFromDateValue invalid date",
			args: args{
				in0: nil,
				in1: nil,
			},
			triggerMock: func() {
				FlagFromDateValue = "any string"
				FlagToDateValue = "2025-03-05"
				FlagTotalDataSampleToGenerateValue = 1000
				FlagPercentageMatchSampleToGenerateValue = 100
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.triggerMock()
			if err := PersistentPreRunner(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunner(t *testing.T) {
	ctx := context.Background()

	type args struct {
		cmd *cobra.Command
		in1 []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				cmd: func() *cobra.Command {
					r := &cobra.Command{}
					r.SetContext(ctx)
					return r
				}(),
				in1: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Runner(tt.args.cmd, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("Runner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

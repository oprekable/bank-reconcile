package helper

import (
	"bytes"
	"context"
	"encoding/csv"
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/default_bank/entity"
)

func TestToBankTrxData(t *testing.T) {
	layoutTime := "2006-01-02 15:04:05"

	type args struct {
		originalData banks.BankTrxDataInterface
		csvReader    *csv.Reader
		filePath     string
		bank         string
		isHaveHeader bool
	}

	tests := []struct {
		name           string
		args           args
		wantReturnData []*banks.BankTrxData
		wantErr        bool
	}{
		{
			name: "Ok with header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: true,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`UniqueIdentifier,Date,Amount
0012d068c53eb0971fc8563343c5d81f,2025-03-15,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14,-42100`,
					)
					return csv.NewReader(f)
				}(),
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "0012d068c53eb0971fc8563343c5d81f",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 00:00:00")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "danamon",
					FilePath: "/foo/bar.csv",
					Amount:   20500,
				},
				{
					UniqueIdentifier: "005dcbc9e27365a072be5393ea8d0f37",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 00:00:00")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "danamon",
					FilePath: "/foo/bar.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error decode with header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: true,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error parse Date",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: true,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`UniqueIdentifier,Date,Amount
0012d068c53eb0971fc8563343c5d81f,random string,20500`,
					)
					return csv.NewReader(f)
				}(),
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Ok without header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: false,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`0012d068c53eb0971fc8563343c5d81f,2025-03-15,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14,-42100`,
					)
					return csv.NewReader(f)
				}(),
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "0012d068c53eb0971fc8563343c5d81f",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 00:00:00")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "danamon",
					FilePath: "/foo/bar.csv",
					Amount:   20500,
				},
				{
					UniqueIdentifier: "005dcbc9e27365a072be5393ea8d0f37",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 00:00:00")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "danamon",
					FilePath: "/foo/bar.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error nil csvReader without header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: false,
				bank:         "danamon",
				csvReader:    nil,
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Error nil originalData without header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: false,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				originalData: nil,
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error decode without header",
			args: args{
				filePath:     "/foo/bar.csv",
				isHaveHeader: false,
				bank:         "danamon",
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				originalData: &entity.CSVBankTrxData{},
			},
			wantReturnData: nil,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnData, err := ToBankTrxData(context.Background(), tt.args.filePath, tt.args.isHaveHeader, tt.args.bank, tt.args.csvReader, tt.args.originalData)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToBankTrxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("ToBankTrxData() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

package bca

import (
	"bytes"
	"context"
	"encoding/csv"
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestBankParserGetBank(t *testing.T) {
	type fields struct {
		bank string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				bank: string(banks.BCABankParser),
			},
			want: string(banks.BCABankParser),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &BankParser{
				bank: tt.fields.bank,
			}

			if got := d.GetBank(); got != tt.want {
				t.Errorf("GetBank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankParserGetParser(t *testing.T) {
	type fields struct {
		parser banks.BankParserType
	}

	tests := []struct {
		name   string
		fields fields
		want   banks.BankParserType
	}{
		{
			name: "Ok",
			fields: fields{
				parser: banks.BCABankParser,
			},
			want: banks.BCABankParser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &BankParser{
				parser: tt.fields.parser,
			}

			if got := d.GetParser(); got != tt.want {
				t.Errorf("GetParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankParserToBankTrxData(t *testing.T) {
	layoutTime := "2006-01-02"

	type fields struct {
		csvReader    *csv.Reader
		parser       banks.BankParserType
		bank         string
		isHaveHeader bool
	}

	type args struct {
		ctx      context.Context
		filePath string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*banks.BankTrxData
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`BCAUniqueIdentifier,BCADate,BCAAmount
0012d068c53eb0971fc8563343c5d81f,2025-03-15,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14,-42100`,
					)
					return csv.NewReader(f)
				}(),
				parser:       banks.BCABankParser,
				bank:         string(banks.BCABankParser),
				isHaveHeader: true,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/foo/bar.csv",
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "0012d068c53eb0971fc8563343c5d81f",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15")
						return t
					}(),
					Type:     banks.CREDIT,
					Bank:     string(banks.BCABankParser),
					FilePath: "/foo/bar.csv",
					Amount:   20500,
				},
				{
					UniqueIdentifier: "005dcbc9e27365a072be5393ea8d0f37",
					Date: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14")
						return t
					}(),
					Type:     banks.DEBIT,
					Bank:     string(banks.BCABankParser),
					FilePath: "/foo/bar.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &BankParser{
				csvReader:    tt.fields.csvReader,
				parser:       tt.fields.parser,
				bank:         tt.fields.bank,
				isHaveHeader: tt.fields.isHaveHeader,
			}

			gotReturnData, err := d.ToBankTrxData(tt.args.ctx, tt.args.filePath)
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

func TestNewBankParser(t *testing.T) {
	type args struct {
		csvReader    *csv.Reader
		bank         string
		isHaveHeader bool
	}

	tests := []struct {
		want    *BankParser
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				bank:         string(banks.BCABankParser),
				csvReader:    csv.NewReader(nil),
				isHaveHeader: false,
			},
			want: &BankParser{
				csvReader:    csv.NewReader(nil),
				parser:       banks.BCABankParser,
				bank:         string(banks.BCABankParser),
				isHaveHeader: false,
			},
			wantErr: false,
		},
		{
			name: "Error nil csvReader",
			args: args{
				bank:         string(banks.BCABankParser),
				csvReader:    nil,
				isHaveHeader: false,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBankParser(tt.args.bank, tt.args.csvReader, tt.args.isHaveHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBankParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBankParser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

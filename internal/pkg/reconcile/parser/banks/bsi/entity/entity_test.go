package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestCSVBankTrxDataGetAbsAmount(t *testing.T) {
	type fields struct {
		BSIAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok - positive",
			fields: fields{
				BSIAmount: 1000,
			},
			want: 1000,
		},
		{
			name: "Ok - negative",
			fields: fields{
				BSIAmount: -1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIAmount: tt.fields.BSIAmount,
			}

			if got := u.GetAbsAmount(); got != tt.want {
				t.Errorf("GetAbsAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetAmount(t *testing.T) {
	type fields struct {
		BSIAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok",
			fields: fields{
				BSIAmount: 1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIAmount: tt.fields.BSIAmount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetBank(t *testing.T) {
	type fields struct {
		BSIBank string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BSIBank: "bsi",
			},
			want: "bsi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIBank: tt.fields.BSIBank,
			}
			if got := u.GetBank(); got != tt.want {
				t.Errorf("GetBank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetDate(t *testing.T) {
	type fields struct {
		BSIDate string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BSIDate: "2025-01-01",
			},
			want: "2025-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIDate: tt.fields.BSIDate,
			}
			if got := u.GetDate(); got != tt.want {
				t.Errorf("GetDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetType(t *testing.T) {
	type fields struct {
		BSIAmount float64
	}
	tests := []struct {
		name   string
		fields fields
		want   banks.TrxType
	}{
		{
			name: "Ok - Debit",
			fields: fields{
				BSIAmount: -1000,
			},
			want: banks.DEBIT,
		},
		{
			name: "Ok - Debit 0",
			fields: fields{
				BSIAmount: 0,
			},
			want: banks.DEBIT,
		},
		{
			name: "Ok - Credit",
			fields: fields{
				BSIAmount: 1000,
			},
			want: banks.CREDIT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIAmount: tt.fields.BSIAmount,
			}
			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetUniqueIdentifier(t *testing.T) {
	type fields struct {
		BSIUniqueIdentifier string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BSIUniqueIdentifier: "bsi-12345",
			},
			want: "bsi-12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIUniqueIdentifier: tt.fields.BSIUniqueIdentifier,
			}
			if got := u.GetUniqueIdentifier(); got != tt.want {
				t.Errorf("GetUniqueIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataToBankTrxData(t *testing.T) {
	layoutTime := "2006-01-02"

	type fields struct {
		BSIUniqueIdentifier string
		BSIDate             string
		BSIBank             string
		BSIAmount           float64
	}
	tests := []struct {
		fields         fields
		wantReturnData *banks.BankTrxData
		name           string
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				BSIUniqueIdentifier: "bsi-12345",
				BSIDate:             "2025-03-15",
				BSIAmount:           20500,
				BSIBank:             "bsi",
			},
			wantReturnData: &banks.BankTrxData{
				UniqueIdentifier: "bsi-12345",
				Date: func() time.Time {
					t, _ := time.Parse(layoutTime, "2025-03-15")
					return t
				}(),
				Type:     banks.CREDIT,
				Bank:     "bsi",
				FilePath: "",
				Amount:   20500,
			},
			wantErr: false,
		},
		{
			name: "Error invalid date format",
			fields: fields{
				BSIUniqueIdentifier: "bsi-12345",
				BSIDate:             "2025/03/15",
				BSIAmount:           20500,
				BSIBank:             "bsi",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BSIUniqueIdentifier: tt.fields.BSIUniqueIdentifier,
				BSIDate:             tt.fields.BSIDate,
				BSIBank:             tt.fields.BSIBank,
				BSIAmount:           tt.fields.BSIAmount,
			}
			gotReturnData, err := u.ToBankTrxData()
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

package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestCSVBankTrxDataGetAbsAmount(t *testing.T) {
	type fields struct {
		BNIAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok - positive",
			fields: fields{
				BNIAmount: 1000,
			},
			want: 1000,
		},
		{
			name: "Ok - negative",
			fields: fields{
				BNIAmount: -1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIAmount: tt.fields.BNIAmount,
			}

			if got := u.GetAbsAmount(); got != tt.want {
				t.Errorf("GetAbsAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetAmount(t *testing.T) {
	type fields struct {
		BNIAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok",
			fields: fields{
				BNIAmount: 1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIAmount: tt.fields.BNIAmount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetBank(t *testing.T) {
	type fields struct {
		BNIBank string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BNIBank: string(banks.BNIBankParser),
			},
			want: string(banks.BNIBankParser),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIBank: tt.fields.BNIBank,
			}

			if got := u.GetBank(); got != tt.want {
				t.Errorf("GetBank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetDate(t *testing.T) {
	type fields struct {
		BNIDate string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BNIDate: "2020-01-01",
			},

			want: "2020-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIDate: tt.fields.BNIDate,
			}

			if got := u.GetDate(); got != tt.want {
				t.Errorf("GetDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetType(t *testing.T) {
	type fields struct {
		BNIAmount float64
	}

	tests := []struct {
		name   string
		want   banks.TrxType
		fields fields
	}{
		{
			name: "DEBIT",
			fields: fields{
				BNIAmount: -1000,
			},
			want: banks.DEBIT,
		},
		{
			name: "CREDIT",
			fields: fields{
				BNIAmount: 1000,
			},
			want: banks.CREDIT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIAmount: tt.fields.BNIAmount,
			}

			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetUniqueIdentifier(t *testing.T) {
	type fields struct {
		BNIUniqueIdentifier string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BNIUniqueIdentifier: "ea270c20-e95b-48ab-b219-b0845fe4e631",
			},
			want: "ea270c20-e95b-48ab-b219-b0845fe4e631",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIUniqueIdentifier: tt.fields.BNIUniqueIdentifier,
			}

			if got := u.GetUniqueIdentifier(); got != tt.want {
				t.Errorf("GetUniqueIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataToBankTrxData(t *testing.T) {
	type fields struct {
		BNIUniqueIdentifier string
		BNIDate             string
		BNIBank             string
		BNIAmount           float64
	}

	tests := []struct {
		wantReturnData *banks.BankTrxData
		name           string
		fields         fields
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				BNIUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				BNIDate:             "1999-01-01",
				BNIBank:             string(banks.BNIBankParser),
				BNIAmount:           1000,
			},
			wantReturnData: &banks.BankTrxData{
				UniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				Date: func() time.Time {
					t, _ := time.Parse("2006-01-02", "1999-01-01")
					return t
				}(),
				Type:     banks.CREDIT,
				Bank:     string(banks.BNIBankParser),
				FilePath: "",
				Amount:   1000,
			},
			wantErr: false,
		},
		{
			name: "Error invalid date",
			fields: fields{
				BNIUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				BNIDate:             "any string",
				BNIBank:             string(banks.BNIBankParser),
				BNIAmount:           1000,
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BNIUniqueIdentifier: tt.fields.BNIUniqueIdentifier,
				BNIDate:             tt.fields.BNIDate,
				BNIBank:             tt.fields.BNIBank,
				BNIAmount:           tt.fields.BNIAmount,
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

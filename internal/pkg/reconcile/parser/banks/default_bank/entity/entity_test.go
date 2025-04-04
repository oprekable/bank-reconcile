package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestCSVBankTrxDataGetAbsAmount(t *testing.T) {
	type fields struct {
		DefaultAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok - positive",
			fields: fields{
				DefaultAmount: 1000,
			},
			want: 1000,
		},
		{
			name: "Ok - negative",
			fields: fields{
				DefaultAmount: -1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultAmount: tt.fields.DefaultAmount,
			}

			if got := u.GetAbsAmount(); got != tt.want {
				t.Errorf("GetAbsAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetAmount(t *testing.T) {
	type fields struct {
		DefaultAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok",
			fields: fields{
				DefaultAmount: 1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultAmount: tt.fields.DefaultAmount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetBank(t *testing.T) {
	type fields struct {
		DefaultBank string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				DefaultBank: "danamon",
			},
			want: "danamon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultBank: tt.fields.DefaultBank,
			}

			if got := u.GetBank(); got != tt.want {
				t.Errorf("GetBank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetDate(t *testing.T) {
	type fields struct {
		DefaultDate string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				DefaultDate: "2020-01-01",
			},

			want: "2020-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultDate: tt.fields.DefaultDate,
			}

			if got := u.GetDate(); got != tt.want {
				t.Errorf("GetDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetType(t *testing.T) {
	type fields struct {
		DefaultAmount float64
	}

	tests := []struct {
		name   string
		want   banks.TrxType
		fields fields
	}{
		{
			name: "DEBIT",
			fields: fields{
				DefaultAmount: -1000,
			},
			want: banks.DEBIT,
		},
		{
			name: "CREDIT",
			fields: fields{
				DefaultAmount: 1000,
			},
			want: banks.CREDIT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultAmount: tt.fields.DefaultAmount,
			}

			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetUniqueIdentifier(t *testing.T) {
	type fields struct {
		DefaultUniqueIdentifier string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				DefaultUniqueIdentifier: "ea270c20-e95b-48ab-b219-b0845fe4e631",
			},
			want: "ea270c20-e95b-48ab-b219-b0845fe4e631",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultUniqueIdentifier: tt.fields.DefaultUniqueIdentifier,
			}

			if got := u.GetUniqueIdentifier(); got != tt.want {
				t.Errorf("GetUniqueIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataToBankTrxData(t *testing.T) {
	type fields struct {
		DefaultUniqueIdentifier string
		DefaultDate             string
		DefaultBank             string
		DefaultAmount           float64
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
				DefaultUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				DefaultDate:             "1999-01-01",
				DefaultBank:             "danamon",
				DefaultAmount:           1000,
			},
			wantReturnData: &banks.BankTrxData{
				UniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				Date: func() time.Time {
					t, _ := time.Parse("2006-01-02", "1999-01-01")
					return t
				}(),
				Type:     banks.CREDIT,
				Bank:     "danamon",
				FilePath: "",
				Amount:   1000,
			},
			wantErr: false,
		},
		{
			name: "Error invalid date",
			fields: fields{
				DefaultUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				DefaultDate:             "any string",
				DefaultBank:             "danamon",
				DefaultAmount:           1000,
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				DefaultUniqueIdentifier: tt.fields.DefaultUniqueIdentifier,
				DefaultDate:             tt.fields.DefaultDate,
				DefaultBank:             tt.fields.DefaultBank,
				DefaultAmount:           tt.fields.DefaultAmount,
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

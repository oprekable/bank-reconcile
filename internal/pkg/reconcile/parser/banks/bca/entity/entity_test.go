package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestCSVBankTrxDataGetAbsAmount(t *testing.T) {
	type fields struct {
		BCAAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok - positive",
			fields: fields{
				BCAAmount: 1000,
			},
			want: 1000,
		},
		{
			name: "Ok - negative",
			fields: fields{
				BCAAmount: -1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCAAmount: tt.fields.BCAAmount,
			}

			if got := u.GetAbsAmount(); got != tt.want {
				t.Errorf("GetAbsAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetAmount(t *testing.T) {
	type fields struct {
		BCAAmount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok",
			fields: fields{
				BCAAmount: 1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCAAmount: tt.fields.BCAAmount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetBank(t *testing.T) {
	type fields struct {
		BCABank string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BCABank: string(banks.BCABankParser),
			},
			want: string(banks.BCABankParser),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCABank: tt.fields.BCABank,
			}

			if got := u.GetBank(); got != tt.want {
				t.Errorf("GetBank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetDate(t *testing.T) {
	type fields struct {
		BCADate string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BCADate: "2020-01-01",
			},

			want: "2020-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCADate: tt.fields.BCADate,
			}

			if got := u.GetDate(); got != tt.want {
				t.Errorf("GetDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetType(t *testing.T) {
	type fields struct {
		BCAAmount float64
	}

	tests := []struct {
		name   string
		want   banks.TrxType
		fields fields
	}{
		{
			name: "DEBIT",
			fields: fields{
				BCAAmount: -1000,
			},
			want: banks.DEBIT,
		},
		{
			name: "CREDIT",
			fields: fields{
				BCAAmount: 1000,
			},
			want: banks.CREDIT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCAAmount: tt.fields.BCAAmount,
			}

			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataGetUniqueIdentifier(t *testing.T) {
	type fields struct {
		BCAUniqueIdentifier string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				BCAUniqueIdentifier: "ea270c20-e95b-48ab-b219-b0845fe4e631",
			},
			want: "ea270c20-e95b-48ab-b219-b0845fe4e631",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCAUniqueIdentifier: tt.fields.BCAUniqueIdentifier,
			}

			if got := u.GetUniqueIdentifier(); got != tt.want {
				t.Errorf("GetUniqueIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVBankTrxDataToBankTrxData(t *testing.T) {
	type fields struct {
		BCAUniqueIdentifier string
		BCADate             string
		BCABank             string
		BCAAmount           float64
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
				BCAUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				BCADate:             "1999-01-01",
				BCABank:             string(banks.BCABankParser),
				BCAAmount:           1000,
			},
			wantReturnData: &banks.BankTrxData{
				UniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				Date: func() time.Time {
					t, _ := time.Parse("2006-01-02", "1999-01-01")
					return t
				}(),
				Type:     banks.CREDIT,
				Bank:     string(banks.BCABankParser),
				FilePath: "",
				Amount:   1000,
			},
			wantErr: false,
		},
		{
			name: "Error invalid date",
			fields: fields{
				BCAUniqueIdentifier: "610085c5-89d7-470b-8158-b4252a9b429d",
				BCADate:             "any string",
				BCABank:             string(banks.BCABankParser),
				BCAAmount:           1000,
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVBankTrxData{
				BCAUniqueIdentifier: tt.fields.BCAUniqueIdentifier,
				BCADate:             tt.fields.BCADate,
				BCABank:             tt.fields.BCABank,
				BCAAmount:           tt.fields.BCAAmount,
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

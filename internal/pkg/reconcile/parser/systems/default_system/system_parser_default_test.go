package default_system

import (
	"bytes"
	"context"
	"encoding/csv"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
)

func TestCSVSystemTrxDataGetAmount(t *testing.T) {
	type fields struct {
		Amount float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Ok - positive",
			fields: fields{
				Amount: 1000,
			},
			want: 1000,
		},
		{
			name: "Ok - negative",
			fields: fields{
				Amount: -1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				Amount: tt.fields.Amount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetTransactionTime(t *testing.T) {
	type fields struct {
		TransactionTime string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok - any string",
			fields: fields{
				TransactionTime: "any string",
			},
			want: "any string",
		},
		{
			name: "Ok - datetime format",
			fields: fields{
				TransactionTime: "1999-11-10 22:58:56",
			},
			want: "1999-11-10 22:58:56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TransactionTime: tt.fields.TransactionTime,
			}
			if got := u.GetTransactionTime(); got != tt.want {
				t.Errorf("GetTransactionTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetTrxID(t *testing.T) {
	type fields struct {
		TrxID string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				TrxID: "123",
			},
			want: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID: tt.fields.TrxID,
			}

			if got := u.GetTrxID(); got != tt.want {
				t.Errorf("GetTrxID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetType(t *testing.T) {
	type fields struct {
		Type string
	}

	tests := []struct {
		name   string
		fields fields
		want   systems.TrxType
	}{
		{
			name: "Ok",
			fields: fields{
				Type: "DEBIT",
			},
			want: "DEBIT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				Type: tt.fields.Type,
			}

			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataToSystemTrxData(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		wantReturnData *systems.SystemTrxData
		name           string
		fields         fields
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				TrxID:           "8ce39d84-b932-43d7-add5-1cacf672b763",
				TransactionTime: "1999-11-10 22:58:56",
				Type:            "DEBIT",
				Amount:          10000,
			},
			wantReturnData: &systems.SystemTrxData{
				TrxID: "8ce39d84-b932-43d7-add5-1cacf672b763",
				TransactionTime: func() time.Time {
					t, _ := time.Parse("2006-01-02 15:04:05", "1999-11-10 22:58:56")
					return t
				}(),
				Type:     "DEBIT",
				FilePath: "",
				Amount:   10000,
			},
			wantErr: false,
		},
		{
			name: "Error time parse",
			fields: fields{
				TrxID:           "8ce39d84-b932-43d7-add5-1cacf672b763",
				TransactionTime: "adasd",
				Type:            "DEBIT",
				Amount:          10000,
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}

			gotReturnData, err := u.ToSystemTrxData()

			if (err != nil) != tt.wantErr {
				t.Errorf("ToSystemTrxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("ToSystemTrxData() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestNewSystemParser(t *testing.T) {
	type args struct {
		dataStruct   systems.SystemTrxDataInterface
		csvReader    *csv.Reader
		isHaveHeader bool
	}

	tests := []struct {
		want    *SystemParser
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    csv.NewReader(nil),
				isHaveHeader: true,
			},
			want: &SystemParser{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    csv.NewReader(nil),
				parser:       systems.DefaultSystemParser,
				isHaveHeader: true,
				poolSystemTrxData: &sync.Pool{
					New: func() interface{} {
						return &systems.SystemTrxData{}
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error csvReader is nil",
			args: args{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    nil,
				isHaveHeader: true,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSystemParser(tt.args.dataStruct, tt.args.csvReader, tt.args.isHaveHeader)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSystemParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewSystemParser() got = %v, want %v", got, tt.want)
				}
			} else {
				if !reflect.DeepEqual(got.parser, tt.want.parser) ||
					!reflect.DeepEqual(got.csvReader, tt.want.csvReader) ||
					!reflect.DeepEqual(got.dataStruct, tt.want.dataStruct) ||
					!reflect.DeepEqual(got.isHaveHeader, tt.want.isHaveHeader) {
					t.Errorf("NewSystemParser() got = %v, want %v", got, tt.want)
				}

				ptr := got.poolSystemTrxData.Get().(*systems.SystemTrxData)
				got.poolSystemTrxData.Put(ptr)
			}
		})
	}
}

func TestSystemParserToSystemTrxData(t *testing.T) {
	layoutTime := "2006-01-02 15:04:05"

	type fields struct {
		dataStruct   systems.SystemTrxDataInterface
		csvReader    *csv.Reader
		parser       systems.SystemParserType
		isHaveHeader bool
	}

	type args struct {
		filePath string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*systems.SystemTrxData
		wantErr        bool
	}{
		{
			name: "Ok with header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`TrxID,TransactionTime,Type,Amount
0012d068c53eb0971fc8563343c5d81f,2025-03-15 10:51:52,CREDIT,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14 18:29:01,CREDIT,42100`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: true,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "0012d068c53eb0971fc8563343c5d81f",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 10:51:52")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   20500,
				},
				{
					TrxID: "005dcbc9e27365a072be5393ea8d0f37",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 18:29:01")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error decode with header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: true,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error parse TransactionTime",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`0012d068c53eb0971fc8563343c5d81f,random string,CREDIT,20500`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Ok without header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`0012d068c53eb0971fc8563343c5d81f,2025-03-15 10:51:52,CREDIT,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14 18:29:01,CREDIT,42100`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "0012d068c53eb0971fc8563343c5d81f",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 10:51:52")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   20500,
				},
				{
					TrxID: "005dcbc9e27365a072be5393ea8d0f37",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 18:29:01")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error nil csvReader without header",
			fields: fields{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    nil,
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Error nil dataStruct without header",
			fields: fields{
				dataStruct: nil,
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error decode without header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SystemParser{
				dataStruct:   tt.fields.dataStruct,
				csvReader:    tt.fields.csvReader,
				parser:       tt.fields.parser,
				isHaveHeader: tt.fields.isHaveHeader,
				poolSystemTrxData: &sync.Pool{
					New: func() interface{} {
						return &systems.SystemTrxData{}
					},
				},
			}

			gotReturnData, err := d.ToSystemTrxData(context.Background(), tt.args.filePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ToSystemTrxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("ToSystemTrxData() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

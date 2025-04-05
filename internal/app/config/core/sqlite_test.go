package core

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/pkg/driver/sql"
)

func TestSqliteParametersOptions(t *testing.T) {
	type fields struct {
		DBPath      string
		Cache       string
		JournalMode string
		IsEnabled   bool
	}

	type args struct {
		logPrefix string
	}

	tests := []struct {
		wantReturnData sql.DBSqliteOption
		name           string
		args           args
		fields         fields
	}{
		{
			name: "Ok",
			fields: fields{
				DBPath:      ":memory:",
				Cache:       "shared",
				JournalMode: "WAL",
				IsEnabled:   false,
			},
			args: args{
				logPrefix: "foo",
			},
			wantReturnData: sql.DBSqliteOption{
				LogPrefix:   "foo",
				DBPath:      ":memory:",
				Cache:       "shared",
				JournalMode: "WAL",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := &SqliteParameters{
				DBPath:      tt.fields.DBPath,
				Cache:       tt.fields.Cache,
				JournalMode: tt.fields.JournalMode,
				IsEnabled:   tt.fields.IsEnabled,
			}

			gotReturnData := pp.Options(tt.args.logPrefix)

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("Options() = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

package sql

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNewSqliteDatabase(t *testing.T) {
	type args struct {
		logger      zerolog.Logger
		option      DBSqliteOption
		isDoLogging bool
	}

	tests := []struct {
		name     string
		args     args
		wantPing bool
		wantErr  bool
	}{
		{
			name: "Ok - no logging",
			args: args{
				option: DBSqliteOption{
					LogPrefix:   "SQLite",
					DBPath:      ":memory:",
					Cache:       "shared",
					JournalMode: "WAL",
				},
				logger:      zerolog.Logger{},
				isDoLogging: false,
			},
			wantErr: false,
		},
		{
			name: "Ok - with logging",
			args: args{
				option: DBSqliteOption{
					LogPrefix:   "SQLite",
					DBPath:      ":memory:",
					Cache:       "shared",
					JournalMode: "WAL",
				},
				logger:      zerolog.Logger{},
				isDoLogging: true,
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				option: DBSqliteOption{
					LogPrefix:   "SQLite",
					DBPath:      ":memory:",
					Cache:       "shared",
					JournalMode: "WAL",
				},
				logger:      zerolog.Logger{},
				isDoLogging: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSqliteDatabase(tt.args.option, tt.args.logger, tt.args.isDoLogging)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSqliteDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

package csqlite

import (
	"bytes"
	"context"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
)

func TestProviderDBSqlite(t *testing.T) {
	var bf bytes.Buffer
	type args struct {
		config *cconfig.Config
		logger *clogger.Logger
		bBPath DBPath
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Err - nil config or logger",
			args: args{
				config: nil,
				logger: nil,
				bBPath: DBPath{},
			},
			wantErr: true,
		},
		{
			name: "Ok",
			args: args{
				config: &cconfig.Config{
					Data: &config.Data{
						Sqlite: core.Sqlite{
							Write: core.SqliteParameters{
								DBPath:    ":memory:",
								IsEnabled: true,
							},
							Read: core.SqliteParameters{
								DBPath:    ":memory:",
								IsEnabled: true,
							},
							IsEnabled: true,
						},
					},
				},
				logger: clogger.NewLogger(
					context.Background(),
					&bf,
				),
				bBPath: DBPath{
					ReadDBPath:  ":memory:",
					WriteDBPath: ":memory:",
				},
			},
			wantErr: false,
		},
		{
			name: "IsEnabled false",
			args: args{
				config: &cconfig.Config{
					Data: &config.Data{
						Sqlite: core.Sqlite{
							IsEnabled: false,
						},
					},
				},
				logger: clogger.NewLogger(
					context.Background(),
					&bf,
				),
				bBPath: DBPath{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn, err := ProviderDBSqlite(tt.args.config, tt.args.logger, tt.args.bBPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderDBSqlite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if fn != nil {
				fn()
			}

			bf.Reset()
		})
	}
}

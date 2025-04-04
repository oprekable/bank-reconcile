package component

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
)

func TestNewComponents(t *testing.T) {
	type args struct {
		config   *cconfig.Config
		logger   *clogger.Logger
		er       *cerror.Error
		dbsqlite *csqlite.DBSqlite
		fs       *cfs.Fs
	}

	tests := []struct {
		args args
		want *Components
		name string
	}{
		{
			name: "Ok",
			args: args{
				config:   &cconfig.Config{},
				logger:   &clogger.Logger{},
				er:       &cerror.Error{},
				dbsqlite: &csqlite.DBSqlite{},
				fs:       &cfs.Fs{},
			},
			want: &Components{
				Config:   &cconfig.Config{},
				Logger:   &clogger.Logger{},
				Error:    &cerror.Error{},
				DBSqlite: &csqlite.DBSqlite{},
				Fs:       &cfs.Fs{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewComponents(tt.args.config, tt.args.logger, tt.args.er, tt.args.dbsqlite, tt.args.fs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

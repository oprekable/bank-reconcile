package component

import (
	"context"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
)

func TestNewComponents(t *testing.T) {
	ctx := context.Background()
	type args struct {
		config   *cconfig.Config
		logger   *clogger.Logger
		er       *cerror.Error
		dbsqlite *csqlite.DBSqlite
		fs       *cfs.Fs
		profiler *cprofiler.Profiler
		ctx      context.Context
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
				profiler: &cprofiler.Profiler{},
				ctx:      ctx,
			},
			want: &Components{
				Config:   &cconfig.Config{},
				Logger:   &clogger.Logger{},
				Error:    &cerror.Error{},
				DBSqlite: &csqlite.DBSqlite{},
				Fs:       &cfs.Fs{},
				Profiler: &cprofiler.Profiler{},
				Context:  ctx,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewComponents(tt.args.ctx, tt.args.config, tt.args.logger, tt.args.er, tt.args.dbsqlite, tt.args.fs, tt.args.profiler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

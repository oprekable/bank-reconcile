package process

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/repository/process/_mock"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/repository/sample/_mock"
)

func TestProviderSvc(t *testing.T) {
	type args struct {
		comp *component.Components
		repo *repository.Repositories
	}

	tests := []struct {
		args args
		want *Svc
		name string
	}{
		{
			name: "Ok",
			args: args{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
					&cprofiler.Profiler{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			want: ProviderSvc(
				component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
					&cprofiler.Profiler{},
				), repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderSvc(tt.args.comp, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderSvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

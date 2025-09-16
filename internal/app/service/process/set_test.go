package process

import (
	"context"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/repository/process/_mock"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/repository/sample/_mock"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

func TestProviderSvc(t *testing.T) {
	ctx := context.Background()
	type args struct {
		comp           *component.Components
		repo           *repository.Repositories
		parserRegistry *banks.ParserRegistry
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
					ctx,
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
				parserRegistry: nil, // Provide nil for the test
			},
			want: ProviderSvc(
				component.NewComponents(
					ctx,
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
				nil, // Provide nil for the test
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderSvc(tt.args.comp, tt.args.repo, tt.args.parserRegistry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderSvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

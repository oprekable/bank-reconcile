package repository

import (
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/repository/process"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/repository/process/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/repository/sample"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/repository/sample/_mock"
)

func TestNewRepositories(t *testing.T) {
	type args struct {
		repoSample  sample.Repository
		repoProcess process.Repository
	}

	tests := []struct {
		args args
		want *Repositories
		name string
	}{
		{
			name: "Ok",
			args: args{
				repoSample:  mocksample.NewRepository(t),
				repoProcess: mockprocess.NewRepository(t),
			},
			want: &Repositories{
				RepoSample:  mocksample.NewRepository(t),
				RepoProcess: mockprocess.NewRepository(t),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepositories(tt.args.repoSample, tt.args.repoProcess); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

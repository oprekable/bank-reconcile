package repository

import (
	"github.com/oprekable/bank-reconcile/internal/app/repository/process"
	"github.com/oprekable/bank-reconcile/internal/app/repository/sample"
)

type Repositories struct {
	RepoSample  sample.Repository
	RepoProcess process.Repository
}

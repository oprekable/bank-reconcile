package component

import (
	"context"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
)

type Components struct {
	Config   *cconfig.Config
	Logger   *clogger.Logger
	Error    *cerror.Error
	DBSqlite *csqlite.DBSqlite
	Fs       *cfs.Fs
	Profiler *cprofiler.Profiler
	Context  context.Context
}

package component

import (
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/cprofiler"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"

	"github.com/spf13/afero"

	"github.com/google/wire"
)

func NewComponents(config *cconfig.Config, logger *clogger.Logger, er *cerror.Error, dbsqlite *csqlite.DBSqlite, fs *cfs.Fs, profiler *cprofiler.Profiler) *Components {
	return &Components{
		Config:   config,
		Logger:   logger,
		Error:    er,
		DBSqlite: dbsqlite,
		Fs:       fs,
		Profiler: profiler,
	}
}

var Set = wire.NewSet(
	wire.Value(
		cconfig.ConfigPaths([]string{
			"./*.toml",
			"./params/*.toml",
		}),
	),
	wire.InterfaceValue(new(afero.Fs), afero.NewOsFs()),
	cconfig.Set,
	clogger.Set,
	cerror.Set,
	csqlite.Set,
	cfs.Set,
	cprofiler.Set,
	NewComponents,
)

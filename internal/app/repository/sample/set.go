package sample

import (
	"github.com/oprekable/bank-reconcile/internal/app/component"

	"github.com/google/wire"
)

func ProviderDB(comp *component.Components) (*DB, error) {
	return NewDB(comp.DBSqlite.DBRead)
}

var Set = wire.NewSet(
	ProviderDB,
	wire.Bind(new(Repository), new(*DB)),
)

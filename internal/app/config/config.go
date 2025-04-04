package config

import (
	"github.com/oprekable/bank-reconcile/internal/app/config/core"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
)

type Data struct {
	core.App
	core.Sqlite
	reconciliation.Reconciliation
}

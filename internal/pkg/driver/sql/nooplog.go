package sql

import (
	"context"

	sqldblogger "github.com/simukti/sqldb-logger"
)

type noopLogAdapter struct {
}

func NewNoopLog() sqldblogger.Logger {
	return &noopLogAdapter{}
}

func (zl *noopLogAdapter) Log(_ context.Context, _ sqldblogger.Level, _ string, _ map[string]interface{}) {
	// do nothing
}

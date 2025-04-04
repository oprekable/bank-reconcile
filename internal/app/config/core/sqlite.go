package core

import (
	"github.com/oprekable/bank-reconcile/internal/pkg/driver/sql"

	"github.com/ulule/deepcopier"
)

// SqliteParameters ..
type SqliteParameters struct {
	DBPath      string `default:":memory:" mapstructure:"db_path"`
	Cache       string `default:"shared"   mapstructure:"cache"`
	JournalMode string `default:"WAL"      mapstructure:"journal_mode"`
	IsEnabled   bool   `default:"false"    mapstructure:"is_enabled"`
}

func (pp *SqliteParameters) Options(logPrefix string) (returnData sql.DBSqliteOption) {
	_ = deepcopier.Copy(pp).To(&returnData)
	returnData.LogPrefix = logPrefix
	return
}

// Sqlite ..
type Sqlite struct {
	Write       SqliteParameters `mapstructure:"write"`
	Read        SqliteParameters `mapstructure:"read"`
	IsEnabled   bool             `default:"false"      mapstructure:"is_enabled"`
	IsDoLogging bool             `default:"false"      mapstructure:"is_do_logging"`
}

package sql

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	_ "modernc.org/sqlite"
)

// DBSqliteOption options for postgres connection
type DBSqliteOption struct {
	LogPrefix   string `deepcopier:"skip"`
	DBPath      string `deepcopier:"field:DBPath"`
	Cache       string `deepcopier:"field:Cache"`
	JournalMode string `deepcopier:"field:JournalMode"`
}

func NewSqliteDatabase(option DBSqliteOption, logger zerolog.Logger, isDoLogging bool) (db *sql.DB, err error) {
	dsn := fmt.Sprintf(
		"%s?cache=%s&_pragma=journal_mode(%s)",
		option.DBPath,
		option.Cache,
		option.JournalMode,
	)

	loggerAdapter := NewNoopLog()

	if isDoLogging {
		loggerAdapter = zerologadapter.New(logger)
	}

	var dbOtel *sql.DB
	driverName := "sqlite"
	if dbOtel, err = otelsql.Open(
		driverName,
		dsn,
		otelsql.WithAttributes(
			semconv.DBSystemSqlite,
		),
	); err == nil {
		db = sqldblogger.OpenDriver(dsn, dbOtel.Driver(), loggerAdapter)

		// Register DB stats to meter
		err = otelsql.RegisterDBStatsMetrics(
			db,
			otelsql.WithAttributes(
				semconv.DBSystemSqlite,
			),
		)
	}

	return db, err
}

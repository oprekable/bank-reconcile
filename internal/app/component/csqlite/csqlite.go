package csqlite

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	sqlDriver "github.com/oprekable/bank-reconcile/internal/pkg/driver/sql"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
)

type DBSqlite struct {
	DBWrite         *sql.DB
	DBRead          *sql.DB
	dBWriteConnOnce sync.Once
	dBReadConnOnce  sync.Once
}

func NewDBSqlite(config *cconfig.Config, logger *clogger.Logger, readDBPath string, writeDBPath string) (rd *DBSqlite, cleanFunc func(), err error) {
	if config == nil || logger == nil {
		err = errors.New("config or logger could not nil")
		return rd, cleanFunc, err
	}

	if !config.Data.Sqlite.IsEnabled {
		return rd, cleanFunc, err
	}

	rd = &DBSqlite{}
	ctx := logger.GetLogger().With().Str("component", "NewDBSqlite").Ctx(context.Background()).Logger().WithContext(logger.GetCtx())

	if config.Data.Sqlite.Write.IsEnabled {
		rd.dBWriteConnOnce.Do(func() {
			dbParameters := config.Data.Sqlite.Write
			if writeDBPath != "" {
				dbParameters.DBPath = writeDBPath
			}

			rd.DBWrite, err = sqlDriver.NewSqliteDatabase(
				dbParameters.Options("sqlite_write"),
				logger.GetLogger(),
				config.Data.Sqlite.IsDoLogging,
			)

			log.AddErr(ctx, err)
		})
	}

	if config.Data.Sqlite.Read.IsEnabled {
		rd.dBReadConnOnce.Do(func() {
			dbParameters := config.Data.Sqlite.Read
			if readDBPath != "" {
				dbParameters.DBPath = readDBPath
			}

			rd.DBRead, err = sqlDriver.NewSqliteDatabase(
				dbParameters.Options("sqlite_read"),
				logger.GetLogger(),
				config.Data.Sqlite.IsDoLogging,
			)

			log.AddErr(ctx, err)
		})
	}

	cleanFunc = func() {
		if rd.DBRead != nil {
			_ = rd.DBRead.Close()
		}

		if rd.DBWrite != nil {
			_ = rd.DBWrite.Close()
		}
	}

	log.Msg(ctx, "sqlite connection loaded")
	return rd, cleanFunc, err
}

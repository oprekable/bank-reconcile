package sample

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/oprekable/bank-reconcile/internal/app/repository/helper"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"

	"github.com/aaronjan/hunch"
	"github.com/goccy/go-json"
)

type DB struct {
	db      *sql.DB
	stmtMap map[string]*sql.Stmt
}

var _ Repository = (*DB)(nil)

func NewDB(
	db *sql.DB,
) (*DB, error) {
	return &DB{
		db:      db,
		stmtMap: make(map[string]*sql.Stmt),
	}, nil
}

func (d *DB) dropTables(ctx context.Context, tx *sql.Tx) (err error) {
	stmtData := []helper.StmtData{
		{
			Name:  "QueryDropTableBanks",
			Query: QueryDropTableBanks,
		},
		{
			Name:  "QueryDropTableArguments",
			Query: QueryDropTableArguments,
		},
		{
			Name:  "QueryDropTableBaseData",
			Query: QueryDropTableBaseData,
		},
	}

	return helper.ExecTxQueries(ctx, tx, d.stmtMap, stmtData)
}

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	return helper.ExecTxQueries(
		ctx,
		tx,
		d.stmtMap,
		[]helper.StmtData{
			{
				Name:  "QueryCreateTableArguments",
				Query: QueryCreateTableArguments,
				Args: func() []any {
					dateStringFormat := "2006-01-02"
					return []any{
						startDate.Format(dateStringFormat),
						toDate.Format(dateStringFormat),
						strconv.FormatInt(limitTrxData, 10),
						strconv.Itoa(matchPercentage),
					}
				}(),
			},
			{
				Name:  "QueryCreateTableBanks",
				Query: QueryCreateTableBanks,
				Args: func() []any {
					b := new(strings.Builder)
					_ = json.NewEncoder(b).Encode(listBank)

					return []any{
						strings.TrimRight(b.String(), "\n"),
					}
				}(),
			},
			{
				Name:  "QueryCreateTableBaseData",
				Query: QueryCreateTableBaseData,
			},
			{
				Name:  "QueryCreateIndexTableBaseData",
				Query: QueryCreateIndexTableBaseData,
			},
		},
	)
}

func (d *DB) postWith(ctx context.Context, methodName string, extraExec hunch.ExecutableInSequence) (err error) {
	var tx *sql.Tx
	defer func() {
		log.Err(ctx, fmt.Sprintf("[sample.NewDB] Exec %s method in db", methodName), helper.CommitOrRollback(tx, err))
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = d.db.BeginTx(ctx, nil)
			return nil, e
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return tx, d.dropTables(c, tx)
		},
		extraExec,
	)

	return
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time, limitTrxData int64, matchPercentage int) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, d.createTables(c, i.(*sql.Tx), listBank, startDate, toDate, limitTrxData, matchPercentage)
	}

	return d.postWith(
		ctx,
		"Pre",
		extraExec,
	)
}

func (d *DB) GetTrx(ctx context.Context) (returnData []TrxData, err error) {
	defer func() {
		log.Err(ctx, "[sample.NewDB] Exec GetData method in db", err)
	}()

	returnData, err = helper.QueryContext[[]TrxData](
		ctx,
		d.db,
		d.stmtMap,
		helper.StmtData{
			Name:  "QueryGetTrxData",
			Query: QueryGetTrxData,
			Args:  nil,
		},
	)

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, nil
	}

	return d.postWith(
		ctx,
		"Post",
		extraExec,
	)
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}

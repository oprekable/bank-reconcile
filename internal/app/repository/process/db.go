package process

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/goccy/go-json"
	"github.com/oprekable/bank-reconcile/internal/app/repository/helper"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
)

const (
	logFlag = "process.NewDB"
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

func (d *DB) dropTableWith(ctx context.Context, methodName string, extraExec hunch.ExecutableInSequence) (err error) {
	execFn := []hunch.ExecutableInSequence{
		func(c context.Context, i interface{}) (r interface{}, e error) {
			tx := i.(*sql.Tx)
			stmtData := []helper.StmtData{
				{
					Name:  "QueryDropTableArguments",
					Query: QueryDropTableArguments,
				},
				{
					Name:  "QueryDropTableBanks",
					Query: QueryDropTableBanks,
				},
				{
					Name:  "QueryDropTableSystemTrx",
					Query: QueryDropTableSystemTrx,
				},
				{
					Name:  "QueryDropTableBankTrx",
					Query: QueryDropTableBankTrx,
				},
				{
					Name:  "QueryDropTableReconciliationMap",
					Query: QueryDropTableReconciliationMap,
				},
			}

			return tx, helper.ExecTxQueries(ctx, tx, d.stmtMap, stmtData)
		},
		extraExec,
	}

	return helper.TxWith(
		ctx,
		logFlag,
		methodName,
		d.db,
		execFn...,
	)
}

func (d *DB) createTables(ctx context.Context, tx *sql.Tx, listBank []string, startDate time.Time, toDate time.Time) (err error) {
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
				Name:  "QueryCreateTableSystemTrx",
				Query: QueryCreateTableSystemTrx,
			},
			{
				Name:  "QueryCreateTableBankTrx",
				Query: QueryCreateTableBankTrx,
			},
			{
				Name:  "QueryCreateTableReconciliationMap",
				Query: QueryCreateTableReconciliationMap,
			},
		},
	)
}

func (d *DB) Pre(ctx context.Context, listBank []string, startDate time.Time, toDate time.Time) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, d.createTables(c, i.(*sql.Tx), listBank, startDate, toDate)
	}

	return d.dropTableWith(
		ctx,
		"Pre",
		extraExec,
	)
}

func (d *DB) importInterface(ctx context.Context, methodName string, query string, data interface{}) (err error) {
	execFn := []hunch.ExecutableInSequence{
		func(c context.Context, i interface{}) (r interface{}, e error) {
			tx := i.(*sql.Tx)
			b, _ := json.Marshal(data)
			stmtData := []helper.StmtData{
				{
					Name:  methodName,
					Query: query,
					Args: func() []any {
						return []any{
							string(b),
						}
					}(),
				},
			}

			return tx, helper.ExecTxQueries(ctx, tx, d.stmtMap, stmtData)
		},
	}

	return helper.TxWith(
		ctx,
		logFlag,
		methodName,
		d.db,
		execFn...,
	)
}

func (d *DB) ImportSystemTrx(ctx context.Context, data []*systems.SystemTrxData, from, to int) (err error) {
	return d.importInterface(ctx, fmt.Sprintf("ImportSystemTrx : range data (%d - %d)", from, to), QueryInsertTableSystemTrx, data)
}

func (d *DB) ImportBankTrx(ctx context.Context, data []*banks.BankTrxData, from, to int) (err error) {
	return d.importInterface(ctx, fmt.Sprintf("ImportBankTrx : range data (%d - %d)", from, to), QueryInsertTableBankTrx, data)
}

func (d *DB) GenerateReconciliationMap(ctx context.Context, minAmount float64, maxAmount float64) (err error) {
	execFn := []hunch.ExecutableInSequence{
		func(c context.Context, i interface{}) (r interface{}, e error) {
			tx := i.(*sql.Tx)
			stmtData := []helper.StmtData{
				{
					Name:  "QueryInsertTableReconciliationMap",
					Query: QueryInsertTableReconciliationMap,
					Args: func() []any {
						return []any{
							minAmount,
							maxAmount,
						}
					}(),
				},
			}

			return tx, helper.ExecTxQueries(ctx, tx, d.stmtMap, stmtData)
		},
	}

	return helper.TxWith(
		ctx,
		logFlag,
		fmt.Sprintf("GenerateReconciliationMap : range amount (%.0f - %.0f)", minAmount, maxAmount),
		d.db,
		execFn...,
	)
}

func (d *DB) GetReconciliationSummary(ctx context.Context) (returnData ReconciliationSummary, err error) {
	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetReconciliationSummary method from db", err)
	}()

	returnData, err = helper.QueryContext[ReconciliationSummary](
		ctx,
		d.db,
		d.stmtMap,
		helper.StmtData{
			Name:  "QueryGetReconciliationSummary",
			Query: QueryGetReconciliationSummary,
			Args:  nil,
		},
	)

	return
}

func (d *DB) Post(ctx context.Context) (err error) {
	extraExec := func(c context.Context, i interface{}) (interface{}, error) {
		return nil, nil
	}

	return d.dropTableWith(
		ctx,
		"Post",
		extraExec,
	)
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}

func (d *DB) GetMatchedTrx(ctx context.Context) (returnData []MatchedTrx, err error) {
	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetMatchedTrx method from db", err)
	}()

	returnData, err = helper.QueryContext[[]MatchedTrx](
		ctx,
		d.db,
		d.stmtMap,
		helper.StmtData{
			Name:  "QueryGetMatchedTrx",
			Query: QueryGetMatchedTrx,
			Args:  nil,
		},
	)

	return
}

func (d *DB) GetNotMatchedSystemTrx(ctx context.Context) (returnData []NotMatchedSystemTrx, err error) {
	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetNotMatchedSystemTrx method from db", err)
	}()

	returnData, err = helper.QueryContext[[]NotMatchedSystemTrx](
		ctx,
		d.db,
		d.stmtMap,
		helper.StmtData{
			Name:  "QueryGetNotMatchedSystemTrx",
			Query: QueryGetNotMatchedSystemTrx,
			Args:  nil,
		},
	)

	return
}

func (d *DB) GetNotMatchedBankTrx(ctx context.Context) (returnData []NotMatchedBankTrx, err error) {
	defer func() {
		log.Err(ctx, "[process.NewDB] Exec GetNotMatchedBankTrx method from db", err)
	}()

	returnData, err = helper.QueryContext[[]NotMatchedBankTrx](
		ctx,
		d.db,
		d.stmtMap,
		helper.StmtData{
			Name:  "QueryGetNotMatchedBankTrx",
			Query: QueryGetNotMatchedBankTrx,
			Args:  nil,
		},
	)

	return
}

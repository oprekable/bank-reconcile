package helper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
	"github.com/pkg/errors"
)

type StmtData struct {
	Name  string
	Query string
	Args  []any
}

func QueryContext[out any](ctx context.Context, db *sql.DB, stmtData StmtData) (returnData out, err error) {
	if db == nil {
		err = errors.New("db is nil")
		return returnData, err
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return db.PrepareContext(c, stmtData.Query) //nolint:sqlclosecheck
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			stmt := i.(*sql.Stmt)
			defer func() {
				_ = stmt.Close() //nolint:sqlclosecheck
			}()
			return stmt.QueryContext(
				c,
				stmtData.Args...,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			switch reflect.TypeOf(returnData).Kind() {
			case reflect.Slice, reflect.Array:
				{
					return nil, scan.RowsStrict(&returnData, rows)
				}
			default:
				return nil, scan.RowStrict(&returnData, rows)
			}
		},
	)

	return returnData, err
}

func ExecTxQueries(ctx context.Context, tx *sql.Tx, stmtData []StmtData) (err error) {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	var executableInSequence []hunch.ExecutableInSequence
	for k := range stmtData {
		executableInSequence = append(
			executableInSequence,
			func(c context.Context, _ interface{}) (r interface{}, e error) {
				stmt, e := tx.PrepareContext(c, stmtData[k].Query) //nolint:sqlclosecheck
				if e != nil {
					return nil, e
				}
				defer func() {
					_ = stmt.Close() //nolint:sqlclosecheck
				}()

				return stmt.ExecContext( //nolint:sqlclosecheck
					c,
					stmtData[k].Args...,
				)
			},
		)
	}

	_, err = hunch.Waterfall(
		ctx,
		executableInSequence...,
	)

	return err
}

func TxWith(ctx context.Context, logFlag string, methodName string, db *sql.DB, extraExec ...hunch.ExecutableInSequence) (err error) {
	var tx *sql.Tx
	defer func() {
		log.Err(ctx, fmt.Sprintf("[%s] Exec %s method in db", logFlag, methodName), CommitOrRollback(tx, err))
	}()

	execFn := []hunch.ExecutableInSequence{
		func(_ context.Context, _ interface{}) (r interface{}, e error) {
			tx, e = db.BeginTx(ctx, nil)
			return tx, e
		},
	}

	execFn = append(
		execFn,
		extraExec...,
	)

	_, err = hunch.Waterfall(
		ctx,
		execFn...,
	)

	return err
}

func CommitOrRollback(tx *sql.Tx, er error) (err error) {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	if er != nil {
		err = errors.Wrap(tx.Rollback(), er.Error())
	} else {
		err = tx.Commit()
	}

	return err
}

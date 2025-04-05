package helper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"

	"github.com/oprekable/bank-reconcile/internal/app/err/core"

	"github.com/aaronjan/hunch"
	"github.com/blockloop/scan/v2"
	"github.com/pkg/errors"
)

type StmtData struct {
	Name  string
	Query string
	Args  []any
}

func QueryContext[out any](ctx context.Context, db *sql.DB, stmtMap map[string]*sql.Stmt, stmtData StmtData) (returnData out, err error) {
	if db == nil {
		err = errors.New("db is nil")
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (_ interface{}, e error) {
			_, ok := stmtMap[stmtData.Name]
			if !ok {
				stmtMap[stmtData.Name], e = db.PrepareContext(c, stmtData.Query) //nolint:sqlclosecheck
			}

			return
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return stmtMap[stmtData.Name].QueryContext(
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

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.CErrDBConn.Error()
	}

	return
}

func ExecTxQueries(ctx context.Context, tx *sql.Tx, stmtMap map[string]*sql.Stmt, stmtData []StmtData) (err error) {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	var executableInSequence []hunch.ExecutableInSequence
	for k := range stmtData {
		executableInSequence = append(
			executableInSequence,
			func(c context.Context, _ interface{}) (r interface{}, e error) {
				defer func() {
					delete(stmtMap, stmtData[k].Name)
				}()

				if _, ok := stmtMap[stmtData[k].Name]; !ok {
					stmtMap[stmtData[k].Name], e = tx.PrepareContext( //nolint:sqlclosecheck
						c,
						stmtData[k].Query,
					)

					if e != nil {
						return nil, e
					}
				}

				return stmtMap[stmtData[k].Name].ExecContext( //nolint:sqlclosecheck
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

	return
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

	return
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

	return
}

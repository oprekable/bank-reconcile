package helper

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/aaronjan/hunch"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	QuerySelectFooBarFaz = "SELECT Bar, Faz FROM Foo WHERE id=?"
	StringRandom         = "random string"
	TwoBar               = "two Bar"
	TwoFaz               = "two Faz"
	QueryInsertFooBarBaz = "INSERT INTO Foo(Bar, Faz) VALUES(?, ?)"
	OneBar               = "one Bar"
	OneFaz               = "one Faz"
)

type Foo struct {
	Bar string `db:"Bar"`
	Faz string `db:"Faz"`
}

type argsQueryContext struct {
	db       *sql.DB
	stmtMap  map[string]*sql.Stmt
	stmtData StmtData
}
type testCaseQueryContext[out any] struct {
	wantReturnData out
	name           string
	args           argsQueryContext
	wantErr        bool
}

func TestCommitOrRollback(t *testing.T) {
	type dbTx struct {
		db *sql.DB
	}

	type args struct {
		dbTx dbTx
		er   error
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "Commit",
			args: args{
				dbTx: func() dbTx {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					s.ExpectCommit()
					return dbTx{
						db: db,
					}
				}(),
				er: nil,
			},
			wantErr: false,
		},
		{
			name: "Rollback",
			args: args{
				dbTx: func() dbTx {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					s.ExpectRollback()
					return dbTx{
						db: db,
					}
				}(),
				er: sql.ErrNoRows,
			},
			wantErr: false,
		},
		{
			name: "Nil tx",
			args: args{
				dbTx: func() dbTx {
					return dbTx{
						db: nil,
					}
				}(),
				er: sql.ErrNoRows,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tx *sql.Tx
			if tt.args.dbTx.db != nil {
				tx, _ = tt.args.dbTx.db.BeginTx(context.Background(), nil)
			}

			if err := CommitOrRollback(tx, tt.args.er); (err != nil) != tt.wantErr {
				t.Errorf("CommitOrRollback() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Cleanup(func() {
				if tt.args.dbTx.db != nil {
					_ = tt.args.dbTx.db.Close()
				}
			})
		})
	}
}

func TestExecTxQueries(t *testing.T) {
	type args struct {
		db       *sql.DB
		stmtMap  map[string]*sql.Stmt
		stmtData []StmtData
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					s.ExpectPrepare(QueryInsertFooBarBaz).
						ExpectExec().
						WithArgs(
							OneBar,
							OneFaz,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: []StmtData{
					{
						Name:  "InsertFoo",
						Query: QueryInsertFooBarBaz,
						Args: []any{
							OneBar,
							OneFaz,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error - PrepareContext",
			args: args{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					s.ExpectPrepare(QueryInsertFooBarBaz).
						WillReturnError(sql.ErrConnDone)
					s.ExpectRollback()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: []StmtData{
					{
						Name:  "InsertFoo",
						Query: QueryInsertFooBarBaz,
						Args: []any{
							OneBar,
							OneFaz,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error - tx nil",
			args: args{
				db:       nil,
				stmtMap:  make(map[string]*sql.Stmt),
				stmtData: []StmtData{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tx *sql.Tx
			if tt.args.db != nil {
				tx, _ = tt.args.db.BeginTx(context.Background(), nil)
			}

			if err := ExecTxQueries(context.Background(), tx, tt.args.stmtMap, tt.args.stmtData); (err != nil) != tt.wantErr {
				t.Errorf("ExecTxQueries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func runQueryContextTest[out any](t *testing.T, tt testCaseQueryContext[out]) {
	gotReturnData, err := QueryContext[out](context.Background(), tt.args.db, tt.args.stmtMap, tt.args.stmtData)
	t.Cleanup(func() {
		if tt.args.db != nil {
			_ = tt.args.db.Close()
		}
	})

	if (err != nil) != tt.wantErr {
		t.Errorf("QueryContext() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
		t.Errorf("QueryContext() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
	}
}

func TestQueryContext(t *testing.T) {

	testsSingleRow := []testCaseQueryContext[Foo]{
		{
			name: "Ok - single row",
			args: argsQueryContext{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QuerySelectFooBarFaz).ExpectQuery().
						WithArgs(StringRandom).
						WillReturnRows(
							sqlmock.NewRows([]string{"Bar", "Faz"}).
								AddRow(OneBar, OneFaz).
								AddRow(TwoBar, TwoFaz))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: QuerySelectFooBarFaz,
					Args:  []any{StringRandom},
				},
			},
			wantReturnData: Foo{
				Bar: OneBar,
				Faz: OneFaz,
			},
			wantErr: false,
		},
		{
			name: "Error sql.ErrNoRows - single row",
			args: argsQueryContext{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QuerySelectFooBarFaz).ExpectQuery().
						WithArgs(StringRandom).
						WillReturnError(sql.ErrNoRows)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: QuerySelectFooBarFaz,
					Args:  []any{StringRandom},
				},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
		{
			name: "Error sql.ErrConnDone - single row",
			args: argsQueryContext{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QuerySelectFooBarFaz).ExpectQuery().
						WithArgs(StringRandom).
						WillReturnError(sql.ErrConnDone)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: QuerySelectFooBarFaz,
					Args:  []any{StringRandom},
				},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
		{
			name: "Error db nil",
			args: argsQueryContext{
				db:       nil,
				stmtMap:  make(map[string]*sql.Stmt),
				stmtData: StmtData{},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
	}

	testsMultipleRow := []testCaseQueryContext[[]Foo]{
		{
			name: "Ok - multiple rows",
			args: argsQueryContext{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QuerySelectFooBarFaz).ExpectQuery().
						WithArgs(StringRandom).
						WillReturnRows(
							sqlmock.NewRows([]string{"Bar", "Faz"}).
								AddRow(OneBar, OneFaz).
								AddRow(TwoBar, TwoFaz))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: QuerySelectFooBarFaz,
					Args:  []any{StringRandom},
				},
			},
			wantReturnData: []Foo{
				{
					Bar: OneBar,
					Faz: OneFaz,
				},
				{
					Bar: TwoBar,
					Faz: TwoFaz,
				},
			},
			wantErr: false,
		},
		{
			name: "Error sql.ErrNoRows - multiple rows",
			args: argsQueryContext{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QuerySelectFooBarFaz).ExpectQuery().
						WithArgs(StringRandom).
						WillReturnError(sql.ErrNoRows)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: QuerySelectFooBarFaz,
					Args:  []any{StringRandom},
				},
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range testsSingleRow {
		t.Run(tt.name, func(t *testing.T) {
			runQueryContextTest[Foo](t, tt)
		})
	}

	for _, tt := range testsMultipleRow {
		t.Run(tt.name, func(t *testing.T) {
			runQueryContextTest[[]Foo](t, tt)
		})
	}
}

func TestTxWith(t *testing.T) {
	type args struct {
		logFlag    string
		methodName string
		db         *sql.DB
		extraExec  []hunch.ExecutableInSequence
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				logFlag:    "Foo.Bar",
				methodName: "FooBar",
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					return db
				}(),
				extraExec: []hunch.ExecutableInSequence{
					func(c context.Context, i interface{}) (r interface{}, e error) {
						tx := i.(*sql.Tx)
						return tx, nil
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TxWith(context.Background(), tt.args.logFlag, tt.args.methodName, tt.args.db, tt.args.extraExec...); (err != nil) != tt.wantErr {
				t.Errorf("TxWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package helper

import (
	"context"
	"database/sql"
	"github.com/aaronjan/hunch"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type Foo struct {
	Bar string `db:"Bar"`
	Faz string `db:"Faz"`
}

func TestCommitOrRollback(t *testing.T) {
	type dbTx struct {
		db *sql.DB
		tx *sql.Tx
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
					r, _ := db.BeginTx(context.Background(), nil)
					s.ExpectCommit()
					return dbTx{
						db: db,
						tx: r,
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
					r, _ := db.BeginTx(context.Background(), nil)
					s.ExpectRollback()
					return dbTx{
						db: db,
						tx: r,
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
						tx: nil,
					}
				}(),
				er: sql.ErrNoRows,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CommitOrRollback(tt.args.dbTx.tx, tt.args.er); (err != nil) != tt.wantErr {
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
		ctx      context.Context
		tx       *sql.Tx
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
				ctx: context.Background(),
				tx: func() *sql.Tx {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()
					r, _ := db.BeginTx(context.Background(), nil)
					s.ExpectPrepare("INSERT INTO Foo(Bar, Faz) VALUES(?, ?)").
						ExpectExec().
						WithArgs(
							"one Bar",
							"one Faz",
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return r
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: []StmtData{
					{
						Name:  "InsertFoo",
						Query: "INSERT INTO Foo(Bar, Faz) VALUES(?, ?)",
						Args: []any{
							"one Bar",
							"one Faz",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error - PrepareContext",
			args: args{
				ctx: context.Background(),
				tx: func() *sql.Tx {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					r, _ := db.BeginTx(context.Background(), nil)
					s.ExpectPrepare("INSERT INTO Foo(Bar, Faz) VALUES(?, ?)").
						WillReturnError(sql.ErrConnDone)
					s.ExpectRollback()

					return r
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: []StmtData{
					{
						Name:  "InsertFoo",
						Query: "INSERT INTO Foo(Bar, Faz) VALUES(?, ?)",
						Args: []any{
							"one Bar",
							"one Faz",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error - tx nil",
			args: args{
				ctx:      context.Background(),
				tx:       nil,
				stmtMap:  make(map[string]*sql.Stmt),
				stmtData: []StmtData{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecTxQueries(tt.args.ctx, tt.args.tx, tt.args.stmtMap, tt.args.stmtData); (err != nil) != tt.wantErr {
				t.Errorf("ExecTxQueries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryContext(t *testing.T) {
	type args struct {
		ctx      context.Context
		db       *sql.DB
		stmtMap  map[string]*sql.Stmt
		stmtData StmtData
	}

	type testCase[out any] struct {
		wantReturnData out
		name           string
		args           args
		wantErr        bool
	}

	testsSingleRow := []testCase[Foo]{
		{
			name: "Ok - single row",
			args: args{
				ctx: context.Background(),
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare("SELECT Bar, Faz FROM Foo WHERE id=?").ExpectQuery().
						WithArgs("random string").
						WillReturnRows(
							sqlmock.NewRows([]string{"Bar", "Faz"}).
								AddRow("one Bar", "one Faz").
								AddRow("two Bar", "two Faz"))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: "SELECT Bar, Faz FROM Foo WHERE id=?",
					Args:  []any{"random string"},
				},
			},
			wantReturnData: Foo{
				Bar: "one Bar",
				Faz: "one Faz",
			},
			wantErr: false,
		},
		{
			name: "Error sql.ErrNoRows - single row",
			args: args{
				ctx: context.Background(),
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare("SELECT Bar, Faz FROM Foo WHERE id=?").ExpectQuery().
						WithArgs("random string").
						WillReturnError(sql.ErrNoRows)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: "SELECT Bar, Faz FROM Foo WHERE id=?",
					Args:  []any{"random string"},
				},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
		{
			name: "Error sql.ErrConnDone - single row",
			args: args{
				ctx: context.Background(),
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare("SELECT Bar, Faz FROM Foo WHERE id=?").ExpectQuery().
						WithArgs("random string").
						WillReturnError(sql.ErrConnDone)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: "SELECT Bar, Faz FROM Foo WHERE id=?",
					Args:  []any{"random string"},
				},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
		{
			name: "Error db nil",
			args: args{
				ctx:      context.Background(),
				db:       nil,
				stmtMap:  make(map[string]*sql.Stmt),
				stmtData: StmtData{},
			},
			wantReturnData: Foo{},
			wantErr:        true,
		},
	}

	testsMultipleRow := []testCase[[]Foo]{
		{
			name: "Ok - multiple rows",
			args: args{
				ctx: context.Background(),
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare("SELECT Bar, Faz FROM Foo WHERE id=?").ExpectQuery().
						WithArgs("random string").
						WillReturnRows(
							sqlmock.NewRows([]string{"Bar", "Faz"}).
								AddRow("one Bar", "one Faz").
								AddRow("two Bar", "two Faz"))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: "SELECT Bar, Faz FROM Foo WHERE id=?",
					Args:  []any{"random string"},
				},
			},
			wantReturnData: []Foo{
				{
					Bar: "one Bar",
					Faz: "one Faz",
				},
				{
					Bar: "two Bar",
					Faz: "two Faz",
				},
			},
			wantErr: false,
		},
		{
			name: "Error sql.ErrNoRows - multiple rows",
			args: args{
				ctx: context.Background(),
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare("SELECT Bar, Faz FROM Foo WHERE id=?").ExpectQuery().
						WithArgs("random string").
						WillReturnError(sql.ErrNoRows)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
				stmtData: StmtData{
					Name:  "SelectFoo",
					Query: "SELECT Bar, Faz FROM Foo WHERE id=?",
					Args:  []any{"random string"},
				},
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range testsSingleRow {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnData, err := QueryContext[Foo](tt.args.ctx, tt.args.db, tt.args.stmtMap, tt.args.stmtData)
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
		})
	}

	for _, tt := range testsMultipleRow {
		t.Run(tt.name, func(t *testing.T) {
			gotReturnData, err := QueryContext[[]Foo](tt.args.ctx, tt.args.db, tt.args.stmtMap, tt.args.stmtData)
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
		})
	}
}

func TestTxWith(t *testing.T) {
	type args struct {
		ctx        context.Context
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
				ctx:        context.Background(),
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
			if err := TxWith(tt.args.ctx, tt.args.logFlag, tt.args.methodName, tt.args.db, tt.args.extraExec...); (err != nil) != tt.wantErr {
				t.Errorf("TxWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

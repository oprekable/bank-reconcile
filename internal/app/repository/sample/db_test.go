package sample

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaronjan/hunch"
)

func TestDBClose(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectClose().WillReturnError(nil)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectClose().WillReturnError(sql.ErrConnDone)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBGetTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []TrxData
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QueryGetTrxData).ExpectQuery().
						WillReturnRows(
							sqlmock.NewRows([]string{
								"trxID",
								"uniqueIdentifier",
								"type",
								"bank",
								"transactionTime",
								"date",
								"is_system_trx",
								"is_bank_trx",
								"amount",
							}).
								AddRow(
									"0012d068c53eb0971fc8563343c5d81f",
									"foo-0012d068c53eb0971fc8563343c5d81f",
									"DEBIT",
									"foo",
									"2025-03-15 10:51:52",
									"2025-03-15",
									true,
									false,
									float64(20500),
								).
								AddRow(
									"005dcbc9e27365a072be5393ea8d0f37",
									"foo-005dcbc9e27365a072be5393ea8d0f37",
									"CREDIT",
									"foo",
									"2025-03-14 18:29:01",
									"2025-03-14",
									true,
									true,
									float64(42100),
								),
						)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx: context.Background(),
			},
			wantReturnData: []TrxData{
				{
					TrxID:            "0012d068c53eb0971fc8563343c5d81f",
					UniqueIdentifier: "foo-0012d068c53eb0971fc8563343c5d81f",
					Type:             "DEBIT",
					Bank:             "foo",
					TransactionTime:  "2025-03-15 10:51:52",
					Date:             "2025-03-15",
					IsSystemTrx:      true,
					IsBankTrx:        false,
					Amount:           20500,
				},
				{
					TrxID:            "005dcbc9e27365a072be5393ea8d0f37",
					UniqueIdentifier: "foo-005dcbc9e27365a072be5393ea8d0f37",
					Type:             "CREDIT",
					Bank:             "foo",
					TransactionTime:  "2025-03-14 18:29:01",
					Date:             "2025-03-14",
					IsSystemTrx:      true,
					IsBankTrx:        true,
					Amount:           42100,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			gotReturnData, err := d.GetTrx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GetTrx() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestDBPost(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		fields  fields
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			if err := d.Post(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBPre(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		startDate       time.Time
		toDate          time.Time
		ctx             context.Context
		listBank        []string
		limitTrxData    int64
		matchPercentage int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableArguments).
						ExpectExec().
						WithArgs(
							"2025-02-28",
							"2025-02-27",
							"1000",
							"100",
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBanks).
						ExpectExec().
						WithArgs(
							`["foo","bar"]`,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateIndexTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx: context.Background(),
				listBank: []string{
					"foo",
					"bar",
				},
				startDate: func() time.Time {
					r, _ := time.Parse("2006-01-02", "2025-02-28")
					return r
				}(),
				toDate: func() time.Time {
					r, _ := time.Parse("2006-01-02", "2025-02-27")
					return r
				}(),
				limitTrxData:    1000,
				matchPercentage: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			if err := d.Pre(tt.args.ctx, tt.args.listBank, tt.args.startDate, tt.args.toDate, tt.args.limitTrxData, tt.args.matchPercentage); (err != nil) != tt.wantErr {
				t.Errorf("Pre() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBCreateTables(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		startDate       time.Time
		toDate          time.Time
		ctx             context.Context
		listBank        []string
		limitTrxData    int64
		matchPercentage int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryCreateTableArguments).
						ExpectExec().
						WithArgs(
							"2025-02-28",
							"2025-02-27",
							"1000",
							"100",
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBanks).
						ExpectExec().
						WithArgs(
							`["foo","bar"]`,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateIndexTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx: context.Background(),
				listBank: []string{
					"foo",
					"bar",
				},
				startDate: func() time.Time {
					r, _ := time.Parse("2006-01-02", "2025-02-28")
					return r
				}(),
				toDate: func() time.Time {
					r, _ := time.Parse("2006-01-02", "2025-02-27")
					return r
				}(),
				limitTrxData:    1000,
				matchPercentage: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			tx, _ := tt.fields.db.BeginTx(tt.args.ctx, nil)
			if err := d.createTables(tt.args.ctx, tx, tt.args.listBank, tt.args.startDate, tt.args.toDate, tt.args.limitTrxData, tt.args.matchPercentage); (err != nil) != tt.wantErr {
				t.Errorf("createTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBDropTables(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		fields  fields
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			tx, _ := tt.fields.db.BeginTx(tt.args.ctx, nil)
			if err := d.dropTables(tt.args.ctx, tx); (err != nil) != tt.wantErr {
				t.Errorf("dropTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBPostWith(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		ctx        context.Context
		extraExec  hunch.ExecutableInSequence
		methodName string
	}

	tests := []struct {
		args    args
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBaseData).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				ctx:        context.Background(),
				methodName: "",
				extraExec: func(c context.Context, i interface{}) (interface{}, error) {
					return nil, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				db:      tt.fields.db,
				stmtMap: tt.fields.stmtMap,
			}

			if err := d.postWith(tt.args.ctx, tt.args.methodName, tt.args.extraExec); (err != nil) != tt.wantErr {
				t.Errorf("postWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDB(t *testing.T) {
	type args struct {
		db *sql.DB
	}

	tests := []struct {
		args    args
		want    *DB
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				db: &sql.DB{},
			},
			want: &DB{
				db:      &sql.DB{},
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDB(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

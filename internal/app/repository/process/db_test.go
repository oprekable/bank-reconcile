package process

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaronjan/hunch"
	"github.com/goccy/go-json"
)

type Foo struct {
	Bar string `db:"Bar"`
	Faz string `db:"Faz"`
}

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

func TestDBGenerateReconciliationMap(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		minAmount float64
		maxAmount float64
	}

	tests := []struct {
		fields  fields
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryInsertTableReconciliationMap).
						ExpectExec().
						WithArgs(float64(0),
							float64(1000),
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				minAmount: 0,
				maxAmount: 1000,
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

			if err := d.GenerateReconciliationMap(context.Background(), tt.args.minAmount, tt.args.maxAmount); (err != nil) != tt.wantErr {
				t.Errorf("GenerateReconciliationMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBGetMatchedTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	tests := []struct {
		name           string
		fields         fields
		wantReturnData []MatchedTrx
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QueryGetMatchedTrx).ExpectQuery().
						WillReturnRows(
							sqlmock.NewRows([]string{"SystemTrxTrxID", "BankTrxUniqueIdentifier", "SystemTrxTransactionTime", "BankTrxDate", "SystemTrxType", "Bank", "SystemTrxAmount", "BankTrxAmount"}).
								AddRow("0012d068c53eb0971fc8563343c5d81f", "foo-0012d068c53eb0971fc8563343c5d81f", "2025-03-15 10:51:52", "2025-03-15", "DEBIT", "foo", 20500, 20500).
								AddRow("005dcbc9e27365a072be5393ea8d0f37", "foo-005dcbc9e27365a072be5393ea8d0f37", "2025-03-14 18:29:01", "2025-03-14", "CREDIT", "foo", 42100, -42100))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantReturnData: []MatchedTrx{
				{
					SystemTrxTrxID:           "0012d068c53eb0971fc8563343c5d81f",
					BankTrxUniqueIdentifier:  "foo-0012d068c53eb0971fc8563343c5d81f",
					SystemTrxTransactionTime: "2025-03-15 10:51:52",
					BankTrxDate:              "2025-03-15",
					SystemTrxType:            "DEBIT",
					Bank:                     "foo",
					SystemTrxAmount:          20500,
					BankTrxAmount:            20500,
				},
				{
					SystemTrxTrxID:           "005dcbc9e27365a072be5393ea8d0f37",
					BankTrxUniqueIdentifier:  "foo-005dcbc9e27365a072be5393ea8d0f37",
					SystemTrxTransactionTime: "2025-03-14 18:29:01",
					BankTrxDate:              "2025-03-14",
					SystemTrxType:            "CREDIT",
					Bank:                     "foo",
					SystemTrxAmount:          42100,
					BankTrxAmount:            -42100,
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

			gotReturnData, err := d.GetMatchedTrx(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchedTrx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GetMatchedTrx() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestDBGetNotMatchedBankTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	tests := []struct {
		name           string
		fields         fields
		wantReturnData []NotMatchedBankTrx
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QueryGetNotMatchedBankTrx).ExpectQuery().
						WillReturnRows(
							sqlmock.NewRows([]string{"UniqueIdentifier", "Date", "Bank", "Amount"}).
								AddRow("0012d068c53eb0971fc8563343c5d81f", "2025-03-15", "foo", 20500).
								AddRow("005dcbc9e27365a072be5393ea8d0f37", "2025-03-14", "foo", 42100))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantReturnData: []NotMatchedBankTrx{
				{
					UniqueIdentifier: "0012d068c53eb0971fc8563343c5d81f",
					Date:             "2025-03-15",
					Bank:             "foo",
					Amount:           20500,
				},
				{
					UniqueIdentifier: "005dcbc9e27365a072be5393ea8d0f37",
					Date:             "2025-03-14",
					Bank:             "foo",
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

			gotReturnData, err := d.GetNotMatchedBankTrx(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotMatchedBankTrx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GetNotMatchedBankTrx() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestDBGetNotMatchedSystemTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	tests := []struct {
		name           string
		fields         fields
		wantReturnData []NotMatchedSystemTrx
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QueryGetNotMatchedSystemTrx).ExpectQuery().
						WillReturnRows(
							sqlmock.NewRows([]string{"TrxID", "TransactionTime", "Type", "Amount"}).
								AddRow("0012d068c53eb0971fc8563343c5d81f", "2025-03-15 10:51:52", "CREDIT", 20500).
								AddRow("005dcbc9e27365a072be5393ea8d0f37", "2025-03-14 18:29:01", "CREDIT", 42100))
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantReturnData: []NotMatchedSystemTrx{
				{
					TrxID:           "0012d068c53eb0971fc8563343c5d81f",
					TransactionTime: "2025-03-15 10:51:52",
					Type:            "CREDIT",
					Amount:          20500,
				},
				{
					TrxID:           "005dcbc9e27365a072be5393ea8d0f37",
					TransactionTime: "2025-03-14 18:29:01",
					Type:            "CREDIT",
					Amount:          42100,
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

			gotReturnData, err := d.GetNotMatchedSystemTrx(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotMatchedSystemTrx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GetNotMatchedSystemTrx() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestDBGetReconciliationSummary(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	tests := []struct {
		fields         fields
		name           string
		wantReturnData ReconciliationSummary
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectPrepare(QueryGetReconciliationSummary).ExpectQuery().
						WillReturnRows(
							sqlmock.NewRows(
								[]string{
									"total_system_trx",
									"total_matched_trx",
									"total_not_matched_trx",
									"sum_system_trx",
									"sum_matched_trx",
									"sum_discrepancies_trx",
								},
							).
								AddRow(1, 1, 0, 1, 1, 0),
						)
					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantReturnData: ReconciliationSummary{
				TotalSystemTrx:      1,
				TotalMatchedTrx:     1,
				TotalNotMatchedTrx:  0,
				SumSystemTrx:        1,
				SumMatchedTrx:       1,
				SumDiscrepanciesTrx: 0,
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

			gotReturnData, err := d.GetReconciliationSummary(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReconciliationSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GetReconciliationSummary() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestDBImportBankTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		data []*banks.BankTrxData
		min  int
		max  int
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
					s.ExpectPrepare(QueryInsertTableBankTrx).
						ExpectExec().
						WithArgs(
							func() string {
								marshal, _ := json.Marshal(
									[]*banks.BankTrxData{
										{
											UniqueIdentifier: "163af765-0769-467f-8185-8ee7166a0098",
											Date:             time.Time{},
											Type:             "DEBIT",
											FilePath:         "/foo/bar",
											Amount:           1000,
										},
									},
								)

								return string(marshal)
							}(),
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				data: []*banks.BankTrxData{
					{
						UniqueIdentifier: "163af765-0769-467f-8185-8ee7166a0098",
						Date:             time.Time{},
						Type:             "DEBIT",
						FilePath:         "/foo/bar",
						Amount:           1000,
					},
				},
				min: 0,
				max: 10,
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

			if err := d.ImportBankTrx(context.Background(), tt.args.data, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("ImportBankTrx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBImportSystemTrx(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		data []*systems.SystemTrxData
		min  int
		max  int
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
					s.ExpectPrepare(QueryInsertTableSystemTrx).
						ExpectExec().
						WithArgs(
							func() string {
								marshal, _ := json.Marshal(
									[]*systems.SystemTrxData{
										{
											TrxID:           "163af765-0769-467f-8185-8ee7166a0098",
											TransactionTime: time.Time{},
											Type:            "DEBIT",
											FilePath:        "/foo/bar",
											Amount:          1000,
										},
									},
								)

								return string(marshal)
							}(),
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				data: []*systems.SystemTrxData{
					{
						TrxID:           "163af765-0769-467f-8185-8ee7166a0098",
						TransactionTime: time.Time{},
						Type:            "DEBIT",
						FilePath:        "/foo/bar",
						Amount:          1000,
					},
				},
				min: 0,
				max: 10,
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

			if err := d.ImportSystemTrx(context.Background(), tt.args.data, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("ImportSystemTrx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBPost(t *testing.T) {
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
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableSystemTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBankTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableReconciliationMap).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
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

			if err := d.Post(context.Background()); (err != nil) != tt.wantErr {
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
		startDate time.Time
		toDate    time.Time
		listBank  []string
	}

	tests := []struct {
		fields  fields
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				db: func() *sql.DB {
					db, s, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
					s.ExpectBegin()

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableSystemTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBankTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableReconciliationMap).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableArguments).
						ExpectExec().
						WithArgs(
							"2025-02-28",
							"2025-02-27",
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBanks).
						ExpectExec().
						WithArgs(
							`["foo","bar"]`,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableSystemTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBankTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableReconciliationMap).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
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

			if err := d.Pre(context.Background(), tt.args.listBank, tt.args.startDate, tt.args.toDate); (err != nil) != tt.wantErr {
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
		startDate time.Time
		toDate    time.Time
		listBank  []string
	}

	tests := []struct {
		fields  fields
		name    string
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
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBanks).
						ExpectExec().
						WithArgs(
							`["foo","bar"]`,
						).
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableSystemTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableBankTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryCreateTableReconciliationMap).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
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

			tx, _ := tt.fields.db.BeginTx(context.Background(), nil)
			if err := d.createTables(context.Background(), tx, tt.args.listBank, tt.args.startDate, tt.args.toDate); (err != nil) != tt.wantErr {
				t.Errorf("createTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBDropTableWith(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
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

					s.ExpectPrepare(QueryDropTableArguments).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBanks).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableSystemTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableBankTrx).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectPrepare(QueryDropTableReconciliationMap).
						ExpectExec().
						WillReturnResult(sqlmock.NewResult(1, 1))

					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
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

			if err := d.dropTableWith(context.Background(), tt.args.methodName, tt.args.extraExec); (err != nil) != tt.wantErr {
				t.Errorf("dropTableWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBImportInterface(t *testing.T) {
	type fields struct {
		db      *sql.DB
		stmtMap map[string]*sql.Stmt
	}

	type args struct {
		data       interface{}
		methodName string
		query      string
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
					s.ExpectPrepare(`INSERT INTO Foo(Bar, Faz) 
						SELECT
							json_extract(j.value, '$.Bar') AS Bar
							 , json_extract(j.value, '$.Faz') AS Faz
						FROM json_each(
							 ?
						) AS j;`).
						ExpectExec().
						WithArgs(
							func() string {
								marshal, _ := json.Marshal(
									[]Foo{
										{
											Bar: "one Bar",
											Faz: "one Faz",
										},
										{
											Bar: "two Bar",
											Faz: "two Faz",
										},
									},
								)

								return string(marshal)
							}(),
						).
						WillReturnResult(sqlmock.NewResult(1, 1))
					s.ExpectCommit()

					return db
				}(),
				stmtMap: make(map[string]*sql.Stmt),
			},
			args: args{
				methodName: "InsertFoo",
				query: `INSERT INTO Foo(Bar, Faz) 
						SELECT
							json_extract(j.value, '$.Bar') AS Bar
							 , json_extract(j.value, '$.Faz') AS Faz
						FROM json_each(
							 ?
						) AS j;`,
				data: []Foo{
					{
						Bar: "one Bar",
						Faz: "one Faz",
					},
					{
						Bar: "two Bar",
						Faz: "two Faz",
					},
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

			if err := d.importInterface(context.Background(), tt.args.methodName, tt.args.query, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("importInterface() error = %v, wantErr %v", err, tt.wantErr)
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

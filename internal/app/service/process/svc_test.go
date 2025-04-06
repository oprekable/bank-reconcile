package process

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/cconfig"
	"github.com/oprekable/bank-reconcile/internal/app/component/cerror"
	"github.com/oprekable/bank-reconcile/internal/app/component/cfs"
	"github.com/oprekable/bank-reconcile/internal/app/component/clogger"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/oprekable/bank-reconcile/internal/app/config/reconciliation"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/repository/process"
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/repository/process/_mock"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/repository/sample/_mock"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"

	"github.com/stretchr/testify/mock"
	"go.chromium.org/luci/common/clock/testclock"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
)

type MockOpenPermissionDeniedFs struct {
	afero.MemMapFs
}

func (o *MockOpenPermissionDeniedFs) Open(_ string) (afero.File, error) {
	return nil, fs.ErrPermission
}

type MockRemoveAllPermissionDeniedFs struct {
	afero.MemMapFs
}

func (r *MockRemoveAllPermissionDeniedFs) RemoveAll(_ string) error {
	return fs.ErrPermission
}

func TestNewSvc(t *testing.T) {
	type args struct {
		comp *component.Components
		repo *repository.Repositories
	}

	tests := []struct {
		args args
		want *Svc
		name string
	}{
		{
			name: "Ok",
			args: args{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			want: NewSvc(
				component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSvc(tt.args.comp, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSvcGenerateReconciliation(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx context.Context
		afs afero.Fs
		bar *progressbar.ProgressBar
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData ReconciliationSummary
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
									SystemTRXPath: "/system",
									BankTRXPath:   "/bank",
									NumberWorker:  10,
									FromDate: func() (r time.Time) {
										r, _ = time.Parse("2006-01-02", "2025-03-06")
										return
									}(),
									ToDate: func() (r time.Time) {
										r, _ = time.Parse("2006-01-02", "2025-03-09")
										return
									}(),
									ListBank: []string{"bca", "bni"},
								},
							},
						}
					}(),
					clogger.NewLogger(
						func() context.Context {
							ctx, _ := testclock.UseTime(context.Background(), time.Unix(1742017753, 0))
							return ctx
						}(),
						&bf,
					),
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"Pre",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"ImportSystemTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"ImportBankTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"GenerateReconciliationMap",
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"GetReconciliationSummary",
							mock.Anything,
						).Return(
							process.ReconciliationSummary{
								TotalSystemTrx:      0,
								TotalMatchedTrx:     0,
								TotalNotMatchedTrx:  0,
								SumSystemTrx:        0,
								SumMatchedTrx:       0,
								SumDiscrepanciesTrx: 0,
							},
							nil,
						).Maybe()

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedBankTrx",
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"Post",
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"Close",
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx: func() context.Context {
					ctx, _ := testclock.UseTime(context.Background(), time.Unix(1742017753, 0))
					return ctx
				}(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					systemTrxFile, _ := f.Create("/system/foo1.csv")
					_, _ = systemTrxFile.Write([]byte(
						`TrxID,TransactionTime,Type,Amount
006630c83821fac6bea13b92b480feb2,2025-03-06 17:09:21,DEBIT,41000
0066a6264a3b04ac25bd93eed2cb3c6c,2025-03-07 10:18:29,CREDIT,1000
0066a6264a3b04ac25bd93eed2cb3aaa,2025-03-07 10:18:29,CREDIT,89900
0066a6264a3b04ac25bd93eed2cb3bbb,2025-03-08 10:18:29,CREDIT,9000
`,
					))

					_ = systemTrxFile.Close()

					bankTrxFile, _ := f.Create("/bank/bca/any_string.csv")
					_, _ = bankTrxFile.Write([]byte(
						`BCAUniqueIdentifier,BCADate,BCAAmount
bca-5585fa85a971917b48ea2729bcf7d9fb,2025-03-06,7700
`,
					))

					_ = bankTrxFile.Close()
					bankTrxFile, _ = f.Create("/bank/bni/any_string.csv")

					_, _ = bankTrxFile.Write([]byte(
						`BNIUniqueIdentifier,BNIDate,BNIAmount
bni-5f4b1bdf10332ea307813ce402f3d7d4,2025-03-09,-71200
`,
					))

					_ = bankTrxFile.Close()

					return f
				}(),
				bar: progressbar.NewOptions(100, progressbar.OptionSetWidth(10), progressbar.OptionSetWriter(&bf)),
			},
			wantReturnData: ReconciliationSummary{
				FileMissingBankTrx:              nil,
				FileMissingSystemTrx:            "",
				FileMatchedSystemTrx:            "",
				TotalProcessedSystemTrx:         0,
				TotalMatchedSystemTrx:           0,
				TotalNotMatchedSystemTrx:        0,
				SumAmountProcessedSystemTrx:     0,
				SumAmountMatchedSystemTrx:       0,
				SumAmountDiscrepanciesSystemTrx: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.GenerateReconciliation(tt.args.ctx, tt.args.afs, tt.args.bar)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateReconciliation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("GenerateReconciliation() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}

			bf.Reset()
		})
	}
}

func TestSvcGenerateReconciliationFiles(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx                   context.Context
		reconciliationSummary *ReconciliationSummary
		fs                    afero.Fs
		isDeleteDirectory     bool
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Ok - nil reconciliationSummary",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args:    args{},
			wantErr: false,
		},
		{
			name: "Error - GetMatchedTrx",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							nil,
							errors.New("GetMatchedTrx error"),
						).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx:                   context.Background(),
				reconciliationSummary: &ReconciliationSummary{},
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: false,
			},
			wantErr: true,
		},
		{
			name: "Error - GetNotMatchedSystemTrx",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							[]process.MatchedTrx{
								{
									SystemTrxTrxID:           "006630c83821fac6bea13b92b480feb2",
									BankTrxUniqueIdentifier:  "bca-5585fa85a971917b48ea2729bcf7d9fb",
									SystemTrxTransactionTime: "2025-03-06 17:09:21",
									BankTrxDate:              "2025-03-06",
									SystemTrxType:            "DEBIT",
									Bank:                     "bca",
									SystemTrxAmount:          41000,
									BankTrxAmount:            41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							nil,
							errors.New("GetNotMatchedSystemTrx error"),
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx:                   context.Background(),
				reconciliationSummary: &ReconciliationSummary{},
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: false,
			},
			wantErr: true,
		},
		{
			name: "Error - GetNotMatchedBankTrx",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							[]process.MatchedTrx{
								{
									SystemTrxTrxID:           "006630c83821fac6bea13b92b480feb2",
									BankTrxUniqueIdentifier:  "bca-5585fa85a971917b48ea2729bcf7d9fb",
									SystemTrxTransactionTime: "2025-03-06 17:09:21",
									BankTrxDate:              "2025-03-06",
									SystemTrxType:            "DEBIT",
									Bank:                     "bca",
									SystemTrxAmount:          41000,
									BankTrxAmount:            41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedSystemTrx{
								{
									TrxID:           "006630c83821fac6bea13b92b480feb2",
									TransactionTime: "2025-03-06 17:09:21",
									Type:            "DEBIT",
									Amount:          41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedBankTrx",
							mock.Anything,
						).Return(
							nil,
							errors.New("GetNotMatchedBankTrx error"),
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx:                   context.Background(),
				reconciliationSummary: &ReconciliationSummary{},
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: false,
			},
			wantErr: true,
		},
		{
			name: "Error - DeleteDirectory",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							[]process.MatchedTrx{
								{
									SystemTrxTrxID:           "006630c83821fac6bea13b92b480feb2",
									BankTrxUniqueIdentifier:  "bca-5585fa85a971917b48ea2729bcf7d9fb",
									SystemTrxTransactionTime: "2025-03-06 17:09:21",
									BankTrxDate:              "2025-03-06",
									SystemTrxType:            "DEBIT",
									Bank:                     "bca",
									SystemTrxAmount:          41000,
									BankTrxAmount:            41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedSystemTrx{
								{
									TrxID:           "006630c83821fac6bea13b92b480feb2",
									TransactionTime: "2025-03-06 17:09:21",
									Type:            "DEBIT",
									Amount:          41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedBankTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedBankTrx{
								{
									UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
									Bank:             "bca",
									Date:             "2025-03-06",
									Amount:           41000,
								},
							},
							nil,
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx:                   context.Background(),
				reconciliationSummary: &ReconciliationSummary{},
				fs: func() afero.Fs {
					f := afero.MemMapFs{}
					_ = f.Mkdir("/report/bank/not_matched", os.ModeDir)
					file, _ := f.Create("/report/bank/not_matched/foo.csv")
					_ = file.Close()
					fDenied := MockRemoveAllPermissionDeniedFs{f}
					return &fDenied
				}(),
				isDeleteDirectory: true,
			},
			wantErr: true,
		},
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							[]process.MatchedTrx{
								{
									SystemTrxTrxID:           "006630c83821fac6bea13b92b480feb2",
									BankTrxUniqueIdentifier:  "bca-5585fa85a971917b48ea2729bcf7d9fb",
									SystemTrxTransactionTime: "2025-03-06 17:09:21",
									BankTrxDate:              "2025-03-06",
									SystemTrxType:            "DEBIT",
									Bank:                     "bca",
									SystemTrxAmount:          41000,
									BankTrxAmount:            41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedSystemTrx{
								{
									TrxID:           "006630c83821fac6bea13b92b480feb2",
									TransactionTime: "2025-03-06 17:09:21",
									Type:            "DEBIT",
									Amount:          41000,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedBankTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedBankTrx{
								{
									UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
									Bank:             "bca",
									Date:             "2025-03-06",
									Amount:           41000,
								},
							},
							nil,
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx:                   context.Background(),
				reconciliationSummary: &ReconciliationSummary{},
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			if err := s.generateReconciliationFiles(tt.args.ctx, tt.args.reconciliationSummary, tt.args.fs, tt.args.isDeleteDirectory); (err != nil) != tt.wantErr {
				t.Errorf("generateReconciliationFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSvcGenerateReconciliationSummaryAndFiles(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx               context.Context
		fs                afero.Fs
		isDeleteDirectory bool
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData ReconciliationSummary
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ReportTRXPath: "/report",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)

						m.On(
							"GetReconciliationSummary",
							mock.Anything,
						).Return(
							process.ReconciliationSummary{
								TotalSystemTrx:      0,
								TotalMatchedTrx:     0,
								TotalNotMatchedTrx:  0,
								SumSystemTrx:        0,
								SumMatchedTrx:       0,
								SumDiscrepanciesTrx: 0,
							},
							nil,
						).Maybe()

						m.On(
							"GetMatchedTrx",
							mock.Anything,
						).Return(
							[]process.MatchedTrx{
								{
									SystemTrxTrxID:           "",
									BankTrxUniqueIdentifier:  "",
									SystemTrxTransactionTime: "",
									BankTrxDate:              "",
									SystemTrxType:            "",
									Bank:                     "",
									SystemTrxAmount:          0,
									BankTrxAmount:            0,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedSystemTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedSystemTrx{
								{
									TrxID:           "",
									TransactionTime: "",
									Type:            "",
									Amount:          0,
								},
							},
							nil,
						).Maybe()

						m.On(
							"GetNotMatchedBankTrx",
							mock.Anything,
						).Return(
							[]process.NotMatchedBankTrx{},
							nil,
						).Maybe()

						return m
					}(),
				),
			},
			args: args{
				ctx: func() context.Context {
					ctx, _ := testclock.UseTime(context.Background(), time.Unix(1742017753, 0))
					return ctx
				}(),
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: true,
			},
			wantReturnData: ReconciliationSummary{
				FileMissingBankTrx:              nil,
				FileMissingSystemTrx:            "/report/system/not_matched/not_matched_1742017753.csv",
				FileMatchedSystemTrx:            "/report/system/matched/matched_1742017753.csv",
				TotalProcessedSystemTrx:         0,
				TotalMatchedSystemTrx:           0,
				TotalNotMatchedSystemTrx:        0,
				SumAmountProcessedSystemTrx:     0,
				SumAmountMatchedSystemTrx:       0,
				SumAmountDiscrepanciesSystemTrx: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.generateReconciliationSummaryAndFiles(tt.args.ctx, tt.args.fs, tt.args.isDeleteDirectory)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateReconciliationSummaryAndFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("generateReconciliationSummaryAndFiles() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestSvcImportReconcileBankDataToDB(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx  context.Context
		data []*banks.BankTrxData
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
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 2,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"ImportBankTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(nil).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				data: []*banks.BankTrxData{
					{
						UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
						Date: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-03-06")
							return t
						}(),
						Type:     "CREDIT",
						Bank:     "BCA",
						FilePath: "/random_string/bca/any_string.csv",
						Amount:   7700,
					},
					{
						UniqueIdentifier: "bni-5f4b1bdf10332ea307813ce402f3d7d4",
						Date: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-03-09")
							return t
						}(),
						Type:     "DEBIT",
						Bank:     "BNI",
						FilePath: "/random_string/bni/any_string.csv",
						Amount:   71200,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 2,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"ImportBankTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(errors.New("error")).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				data: []*banks.BankTrxData{
					{
						UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
						Date: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-03-06")
							return t
						}(),
						Type:     "CREDIT",
						Bank:     "BCA",
						FilePath: "/random_string/bca/any_string.csv",
						Amount:   7700,
					},
					{
						UniqueIdentifier: "bni-5f4b1bdf10332ea307813ce402f3d7d4",
						Date: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-03-09")
							return t
						}(),
						Type:     "DEBIT",
						Bank:     "BNI",
						FilePath: "/random_string/bni/any_string.csv",
						Amount:   71200,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			if err := s.importReconcileBankDataToDB(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("importReconcileBankDataToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSvcImportReconcileMapToDB(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx context.Context
		min float64
		max float64
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
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 2,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"GenerateReconciliationMap",
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(nil).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				min: 1,
				max: 10,
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 2,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"GenerateReconciliationMap",
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(errors.New("error")).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				min: 1,
				max: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			if err := s.importReconcileMapToDB(tt.args.ctx, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("importReconcileMapToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSvcImportReconcileSystemDataToDB(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx  context.Context
		data []*systems.SystemTrxData
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
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 10,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"ImportSystemTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(nil).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				data: []*systems.SystemTrxData{
					{
						TrxID: "0066a6264a3b04ac25bd93eed2cb3c6c",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
							return t
						}(),
						Type:     "CREDIT",
						FilePath: "/foo2.csv",
						Amount:   41000,
					},
					{
						TrxID: "006630c83821fac6bea13b92b480feb2",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-11 17:09:21")
							return t
						}(),
						Type:     "DEBIT",
						FilePath: "/foo1.csv",
						Amount:   89900,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									NumberWorker: 10,
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					func() process.Repository {
						m := mockprocess.NewRepository(t)
						m.On(
							"ImportSystemTrx",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(errors.New("error")).Maybe()
						return m
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				data: []*systems.SystemTrxData{
					{
						TrxID: "0066a6264a3b04ac25bd93eed2cb3c6c",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
							return t
						}(),
						Type:     "CREDIT",
						FilePath: "/foo2.csv",
						Amount:   41000,
					},
					{
						TrxID: "006630c83821fac6bea13b92b480feb2",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-11 17:09:21")
							return t
						}(),
						Type:     "DEBIT",
						FilePath: "/foo1.csv",
						Amount:   89900,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			if err := s.importReconcileSystemDataToDB(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("importReconcileSystemDataToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSvcParse(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx context.Context
		afs afero.Fs
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantTrxData parser.TrxData
		wantErr     bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									FromDate: func() time.Time {
										t, _ := time.Parse("2006-01-02", "2025-03-06")
										return t
									}(),
									ToDate: func() time.Time {
										t, _ := time.Parse("2006-01-02", "2025-03-07")
										return t
									}(),
									SystemTRXPath: "/system",
									BankTRXPath:   "/bank",
									ListBank:      []string{"bca", "bni"},
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					systemTrxFile, _ := f.Create("/system/foo1.csv")
					_, _ = systemTrxFile.Write([]byte(
						`TrxID,TransactionTime,Type,Amount
006630c83821fac6bea13b92b480feb2,2025-03-06 17:09:21,DEBIT,41000
0066a6264a3b04ac25bd93eed2cb3c6c,2025-03-07 10:18:29,CREDIT,1000
0066a6264a3b04ac25bd93eed2cb3aaa,2025-03-07 10:18:29,CREDIT,89900
0066a6264a3b04ac25bd93eed2cb3bbb,2025-03-08 10:18:29,CREDIT,9000
`,
					))

					_ = systemTrxFile.Close()

					bankTrxFile, _ := f.Create("/bank/bca/any_string.csv")
					_, _ = bankTrxFile.Write([]byte(
						`BCAUniqueIdentifier,BCADate,BCAAmount
bca-5585fa85a971917b48ea2729bcf7d9fb,2025-03-06,7700
`,
					))

					_ = bankTrxFile.Close()
					bankTrxFile, _ = f.Create("/bank/bni/any_string.csv")

					_, _ = bankTrxFile.Write([]byte(
						`BNIUniqueIdentifier,BNIDate,BNIAmount
bni-5f4b1bdf10332ea307813ce402f3d7d4,2025-03-09,-71200
`,
					))

					_ = bankTrxFile.Close()

					return f
				}(),
			},
			wantTrxData: parser.TrxData{
				SystemTrx: []*systems.SystemTrxData{
					{
						TrxID: "006630c83821fac6bea13b92b480feb2",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-06 17:09:21")
							return t
						}(),
						Type:     "DEBIT",
						FilePath: "/system/foo1.csv",
						Amount:   41000,
					},
					{
						TrxID: "0066a6264a3b04ac25bd93eed2cb3c6c",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
							return t
						}(),
						Type:     "CREDIT",
						FilePath: "/system/foo1.csv",
						Amount:   1000,
					},
					{
						TrxID: "0066a6264a3b04ac25bd93eed2cb3aaa",
						TransactionTime: func() time.Time {
							t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
							return t
						}(),
						Type:     "CREDIT",
						FilePath: "/system/foo1.csv",
						Amount:   89900,
					},
				},
				BankTrx: []*banks.BankTrxData{
					{
						UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
						Date: func() time.Time {
							t, _ := time.Parse("2006-01-02", "2025-03-06")
							return t
						}(),
						Type:     "CREDIT",
						Bank:     "BCA",
						FilePath: "/bank/bca/any_string.csv",
						Amount:   7700,
					},
				},
				MinSystemAmount: 0,
				MaxSystemAmount: 89900,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotTrxData, err := s.parse(tt.args.ctx, tt.args.afs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotTrxData, tt.wantTrxData) {
				t.Errorf("parse() gotTrxData = %v, want %v", gotTrxData, tt.wantTrxData)
			}
		})
	}
}

func TestSvcParseBankTrxFile(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx  context.Context
		afs  afero.Fs
		item FilePathBankTrx
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*banks.BankTrxData
		wantErr        bool
	}{
		{
			name: "Ok bca",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/bca.csv")

					_, _ = fooFile.Write([]byte(
						`BCAUniqueIdentifier,BCADate,BCAAmount
bca-e6f8fbe1f6f8c72da7caade610b692e8,2025-03-04,-71700
bca-5585fa85a971917b48ea2729bcf7d9fb,2025-03-06,7700
`,
					))

					_ = fooFile.Close()
					return f
				}(),
				item: FilePathBankTrx{
					Bank:     "bca",
					FilePath: "/bca.csv",
				},
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "bca-e6f8fbe1f6f8c72da7caade610b692e8",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-04")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "BCA",
					FilePath: "/bca.csv",
					Amount:   71700,
				},
				{
					UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-06")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "BCA",
					FilePath: "/bca.csv",
					Amount:   7700,
				},
			},
			wantErr: false,
		},
		{
			name: "Ok bni",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/bni.csv")

					_, _ = fooFile.Write([]byte(
						`BNIUniqueIdentifier,BNIDate,BNIAmount
bni-7b422b9abac7a628125bc1c6bc7adced,2025-03-04,79500
bni-5f4b1bdf10332ea307813ce402f3d7d4,2025-03-09,-71200
`,
					))

					_ = fooFile.Close()
					return f
				}(),
				item: FilePathBankTrx{
					Bank:     "bni",
					FilePath: "/bni.csv",
				},
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "bni-7b422b9abac7a628125bc1c6bc7adced",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-04")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "BNI",
					FilePath: "/bni.csv",
					Amount:   79500,
				},
				{
					UniqueIdentifier: "bni-5f4b1bdf10332ea307813ce402f3d7d4",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-09")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "BNI",
					FilePath: "/bni.csv",
					Amount:   71200,
				},
			},
			wantErr: false,
		},
		{
			name: "Ok default",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/foo.csv")

					_, _ = fooFile.Write([]byte(
						`UniqueIdentifier,Date,Amount
foo-7b422b9abac7a628125bc1c6bc7adced,2025-03-04,79500
foo-5f4b1bdf10332ea307813ce402f3d7d4,2025-03-09,-71200
`,
					))

					_ = fooFile.Close()
					return f
				}(),
				item: FilePathBankTrx{
					Bank:     "foo",
					FilePath: "/foo.csv",
				},
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "foo-7b422b9abac7a628125bc1c6bc7adced",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-04")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "FOO",
					FilePath: "/foo.csv",
					Amount:   79500,
				},
				{
					UniqueIdentifier: "foo-5f4b1bdf10332ea307813ce402f3d7d4",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-09")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "FOO",
					FilePath: "/foo.csv",
					Amount:   71200,
				},
			},
			wantErr: false,
		},
		{
			name: "Error ToBankTrxData",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/foo.csv")
					_ = f.Chmod("/foo.csv", 0000)

					_, _ = fooFile.Write([]byte(
						`UniqueIdentifier
foo-7b422b9abac7a628125bc1c6bc7adced,string,79500
`,
					))

					_ = fooFile.Close()
					return f
				}(),
				item: FilePathBankTrx{
					Bank:     "foo",
					FilePath: "/foo.csv",
				},
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error File",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.MemMapFs{}
					fDenied := MockOpenPermissionDeniedFs{f}
					fooFile, _ := fDenied.Create("/foo.csv")
					_ = fooFile.Close()
					return &fDenied
				}(),
				item: FilePathBankTrx{
					Bank:     "foo",
					FilePath: "/foo.csv",
				},
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.parseBankTrxFile(tt.args.ctx, tt.args.afs, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBankTrxFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("parseBankTrxFile() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestSvcParseBankTrxFiles(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx context.Context
		afs afero.Fs
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*banks.BankTrxData
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									BankTRXPath: "/random_string",
									ListBank:    []string{"bca", "bni"},
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/random_string/bca/any_string.csv")
					_, _ = fooFile.Write([]byte(
						`BCAUniqueIdentifier,BCADate,BCAAmount
bca-5585fa85a971917b48ea2729bcf7d9fb,2025-03-06,7700
`,
					))

					_ = fooFile.Close()

					fooFile, _ = f.Create("/random_string/bni/any_string.csv")

					_, _ = fooFile.Write([]byte(
						`BNIUniqueIdentifier,BNIDate,BNIAmount
bni-5f4b1bdf10332ea307813ce402f3d7d4,2025-03-09,-71200
`,
					))

					_ = fooFile.Close()
					return f
				}(),
			},
			wantReturnData: []*banks.BankTrxData{
				{
					UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-06")
						return t
					}(),
					Type:     "CREDIT",
					Bank:     "BCA",
					FilePath: "/random_string/bca/any_string.csv",
					Amount:   7700,
				},
				{
					UniqueIdentifier: "bni-5f4b1bdf10332ea307813ce402f3d7d4",
					Date: func() time.Time {
						t, _ := time.Parse("2006-01-02", "2025-03-09")
						return t
					}(),
					Type:     "DEBIT",
					Bank:     "BNI",
					FilePath: "/random_string/bni/any_string.csv",
					Amount:   71200,
				},
			},
			wantErr: false,
		},
		{
			name: "Ok - bank not in the list",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									BankTRXPath: "/random_string",
									ListBank:    []string{"foo"},
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/random_string/bca/any_string.csv")
					_, _ = fooFile.Write([]byte(
						`BCAUniqueIdentifier,BCADate,BCAAmount
bca-5585fa85a971917b48ea2729bcf7d9fb,2025-03-06,7700
`,
					))

					_ = fooFile.Close()
					return f
				}(),
			},
			wantReturnData: nil,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.parseBankTrxFiles(tt.args.ctx, tt.args.afs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBankTrxFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(gotReturnData, func(i, j int) bool {
				return gotReturnData[i].UniqueIdentifier < gotReturnData[j].UniqueIdentifier
			})

			sort.Slice(tt.wantReturnData, func(i, j int) bool {
				return tt.wantReturnData[i].UniqueIdentifier < tt.wantReturnData[j].UniqueIdentifier
			})

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("parseBankTrxFiles() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestSvcParseSystemTrxFile(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx      context.Context
		afs      afero.Fs
		filePath string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*systems.SystemTrxData
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/foo.csv")

					_, _ = fooFile.Write([]byte(
						`TrxID,TransactionTime,Type,Amount
006630c83821fac6bea13b92b480feb2,2025-03-11 17:09:21,DEBIT,89900
0066a6264a3b04ac25bd93eed2cb3c6c,2025-03-07 10:18:29,CREDIT,41000
`,
					))

					_ = fooFile.Close()
					return f
				}(),
				filePath: "/foo.csv",
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "006630c83821fac6bea13b92b480feb2",
					TransactionTime: func() time.Time {
						t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-11 17:09:21")
						return t
					}(),
					Type:     "DEBIT",
					FilePath: "/foo.csv",
					Amount:   89900,
				},
				{
					TrxID: "0066a6264a3b04ac25bd93eed2cb3c6c",
					TransactionTime: func() time.Time {
						t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/foo.csv",
					Amount:   41000,
				},
			},
			wantErr: false,
		},
		{
			name: "Error file not found",
			fields: fields{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				filePath: "/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.parseSystemTrxFile(tt.args.ctx, tt.args.afs, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSystemTrxFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("parseSystemTrxFile() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestSvcParseSystemTrxFiles(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx context.Context
		afs afero.Fs
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*systems.SystemTrxData
		wantErr        bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									SystemTRXPath: "/",
								},
							},
						}
					}(),
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			args: args{
				ctx: context.Background(),
				afs: func() afero.Fs {
					f := afero.NewMemMapFs()
					fooFile, _ := f.Create("/foo1.csv")

					_, _ = fooFile.Write([]byte(
						`TrxID,TransactionTime,Type,Amount
006630c83821fac6bea13b92b480feb2,2025-03-11 17:09:21,DEBIT,89900
`,
					))

					_ = fooFile.Close()
					fooFile, _ = f.Create("/foo2.csv")

					_, _ = fooFile.Write([]byte(
						`TrxID,TransactionTime,Type,Amount
0066a6264a3b04ac25bd93eed2cb3c6c,2025-03-07 10:18:29,CREDIT,41000
`,
					))

					_ = fooFile.Close()
					return f
				}(),
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "0066a6264a3b04ac25bd93eed2cb3c6c",
					TransactionTime: func() time.Time {
						t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-07 10:18:29")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/foo2.csv",
					Amount:   41000,
				},
				{
					TrxID: "006630c83821fac6bea13b92b480feb2",
					TransactionTime: func() time.Time {
						t, _ := time.Parse("2006-01-02 15:04:05", "2025-03-11 17:09:21")
						return t
					}(),
					Type:     "DEBIT",
					FilePath: "/foo1.csv",
					Amount:   89900,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotReturnData, err := s.parseSystemTrxFiles(tt.args.ctx, tt.args.afs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSystemTrxFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(gotReturnData, func(i, j int) bool {
				return gotReturnData[i].TrxID < gotReturnData[j].TrxID
			})

			sort.Slice(tt.wantReturnData, func(i, j int) bool {
				return tt.wantReturnData[i].TrxID < tt.wantReturnData[j].TrxID
			})

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("parseSystemTrxFiles() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

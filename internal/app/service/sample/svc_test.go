package sample

import (
	"bytes"
	"context"
	"reflect"
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
	mockprocess "github.com/oprekable/bank-reconcile/internal/app/repository/process/_mock"
	"github.com/oprekable/bank-reconcile/internal/app/repository/sample"
	mocksample "github.com/oprekable/bank-reconcile/internal/app/repository/sample/_mock"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	entitybca "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bca/entity"
	entitybni "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bni/entity"
	entitydefaultbank "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/default_bank/entity"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems/default_system"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"go.chromium.org/luci/common/clock/testclock"
)

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

func TestSvcGenerateSample(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		ctx               context.Context
		fs                afero.Fs
		bar               *progressbar.ProgressBar
		isDeleteDirectory bool
	}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantReturnSummary Summary
		wantErr           bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									ListBank: []string{"bca"},
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
					func() sample.Repository {
						m := mocksample.NewRepository(t)
						m.On(
							"Pre",
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
							mock.Anything,
						).Return(
							nil,
							nil,
						).Maybe()

						m.On(
							"GetTrx",
							mock.Anything,
						).Return(
							[]sample.TrxData{
								{
									TrxID:            "006630c83821fac6bea13b92b480feb2",
									UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
									Type:             "DEBIT",
									Bank:             "bca",
									TransactionTime:  "2025-03-06 17:09:21",
									Date:             "2025-03-06",
									IsSystemTrx:      true,
									IsBankTrx:        true,
									Amount:           0,
								},
							},
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
					mockprocess.NewRepository(t),
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
				bar:               progressbar.NewOptions(100, progressbar.OptionSetWidth(10), progressbar.OptionSetWriter(&bf)),
				isDeleteDirectory: false,
			},
			wantReturnSummary: Summary{
				TotalBankTrx: func() map[string]int64 {
					m := make(map[string]int64)
					m["bca"] = 1
					return m
				}(),
				FileBankTrx: func() map[string]string {
					m := make(map[string]string)
					m["bca"] = "/bca/bca_1742017753.csv"
					return m
				}(),
				FileSystemTrx:  "/1742017753.csv",
				TotalSystemTrx: 1,
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

			gotReturnSummary, err := s.GenerateSample(tt.args.ctx, tt.args.fs, tt.args.bar, tt.args.isDeleteDirectory)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnSummary, tt.wantReturnSummary) {
				t.Errorf("GenerateSample() gotReturnSummary = %v, want %v", gotReturnSummary, tt.wantReturnSummary)
			}

			bf.Reset()
		})
	}
}

func TestSvcAppendExecutor(t *testing.T) {
	type fields struct {
		comp *component.Components
		repo *repository.Repositories
	}

	type args struct {
		fs                afero.Fs
		filePath          string
		trxDataSlice      []banks.BankTrxDataInterface
		isDeleteDirectory bool
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantTotalData int64
	}{
		{
			name: "Ok - zero data",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				filePath:          "/foo.csv",
				trxDataSlice:      nil,
				isDeleteDirectory: false,
			},
			wantTotalData: 0,
		},
		{
			name: "Ok - bca",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				filePath: "/foo.csv",
				trxDataSlice: func() []banks.BankTrxDataInterface {
					return []banks.BankTrxDataInterface{
						&entitybca.CSVBankTrxData{
							BCAUniqueIdentifier: "",
							BCADate:             "",
							BCABank:             "",
							BCAAmount:           0,
						},
					}
				}(),
				isDeleteDirectory: false,
			},
			wantTotalData: 1,
		},
		{
			name: "Ok - bni",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				filePath: "/foo.csv",
				trxDataSlice: func() []banks.BankTrxDataInterface {
					return []banks.BankTrxDataInterface{
						&entitybni.CSVBankTrxData{
							BNIUniqueIdentifier: "",
							BNIDate:             "",
							BNIBank:             "",
							BNIAmount:           0,
						},
					}
				}(),
				isDeleteDirectory: false,
			},
			wantTotalData: 1,
		},
		{
			name: "Ok - default",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				filePath: "/foo.csv",
				trxDataSlice: func() []banks.BankTrxDataInterface {
					return []banks.BankTrxDataInterface{
						&entitydefaultbank.CSVBankTrxData{
							DefaultUniqueIdentifier: "",
							DefaultDate:             "",
							DefaultBank:             "",
							DefaultAmount:           0,
						},
					}
				}(),
				isDeleteDirectory: false,
			},
			wantTotalData: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotTotalData, gotExecutor := s.appendExecutor(tt.args.fs, tt.args.filePath, tt.args.trxDataSlice, tt.args.isDeleteDirectory)
			if gotTotalData != tt.wantTotalData {
				t.Errorf("appendExecutor() gotTotalData = %v, want %v", gotTotalData, tt.wantTotalData)
			}

			if gotExecutor != nil {
				_, err := gotExecutor(context.Background())
				if err != nil {
					t.Errorf("appendExecutor() error = %v, wantErr no error", err)
					return
				}
			}
		})
	}
}

func TestSvcDeleteDirectorySystemTrxBankTrx(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Ok - isDeleteDirectory = true",
			fields: fields{
				comp: component.NewComponents(
					func() *cconfig.Config {
						return &cconfig.Config{
							Data: &config.Data{
								Reconciliation: reconciliation.Reconciliation{
									SystemTRXPath: "/system",
									BankTRXPath:   "/bank",
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
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: true,
			},
			wantErr: false,
		},
		{
			name: "Ok - isDeleteDirectory = false",
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
				fs: func() afero.Fs {
					f := afero.NewMemMapFs()
					return f
				}(),
				isDeleteDirectory: false,
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

			if err := s.deleteDirectorySystemTrxBankTrx(tt.args.ctx, tt.args.fs, tt.args.isDeleteDirectory); (err != nil) != tt.wantErr {
				t.Errorf("deleteDirectorySystemTrxBankTrx() error = %v, wantErr %v", err, tt.wantErr)
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
		data sample.TrxData
	}

	tests := []struct {
		fields            fields
		wantSystemTrxData systems.SystemTrxDataInterface
		wantBankTrxData   banks.BankTrxDataInterface
		name              string
		args              args
	}{
		{
			name: "Ok - bca",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				data: sample.TrxData{
					TrxID:            "006630c83821fac6bea13b92b480feb2",
					UniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
					Type:             "DEBIT",
					Bank:             "bca",
					TransactionTime:  "2025-03-06 17:09:21",
					Date:             "2025-03-06",
					IsSystemTrx:      true,
					IsBankTrx:        true,
					Amount:           41000,
				},
			},
			wantSystemTrxData: &default_system.CSVSystemTrxData{
				TrxID:           "006630c83821fac6bea13b92b480feb2",
				TransactionTime: "2025-03-06 17:09:21",
				Type:            "DEBIT",
				Amount:          41000,
			},
			wantBankTrxData: &entitybca.CSVBankTrxData{
				BCAUniqueIdentifier: "bca-5585fa85a971917b48ea2729bcf7d9fb",
				BCADate:             "2025-03-06",
				BCAAmount:           -41000,
				BCABank:             "bca",
			},
		},
		{
			name: "Ok - bni",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				data: sample.TrxData{
					TrxID:            "006630c83821fac6bea13b92b480feb2",
					UniqueIdentifier: "bni-5585fa85a971917b48ea2729bcf7d9fb",
					Type:             "DEBIT",
					Bank:             "bni",
					TransactionTime:  "2025-03-06 17:09:21",
					Date:             "2025-03-06",
					IsSystemTrx:      true,
					IsBankTrx:        true,
					Amount:           41000,
				},
			},
			wantSystemTrxData: &default_system.CSVSystemTrxData{
				TrxID:           "006630c83821fac6bea13b92b480feb2",
				TransactionTime: "2025-03-06 17:09:21",
				Type:            "DEBIT",
				Amount:          41000,
			},
			wantBankTrxData: &entitybni.CSVBankTrxData{
				BNIUniqueIdentifier: "bni-5585fa85a971917b48ea2729bcf7d9fb",
				BNIDate:             "2025-03-06",
				BNIAmount:           -41000,
				BNIBank:             "bni",
			},
		},
		{
			name: "Ok - default",
			fields: fields{
				comp: nil,
				repo: nil,
			},
			args: args{
				data: sample.TrxData{
					TrxID:            "006630c83821fac6bea13b92b480feb2",
					UniqueIdentifier: "foo-5585fa85a971917b48ea2729bcf7d9fb",
					Type:             "DEBIT",
					Bank:             "foo",
					TransactionTime:  "2025-03-06 17:09:21",
					Date:             "2025-03-06",
					IsSystemTrx:      true,
					IsBankTrx:        true,
					Amount:           41000,
				},
			},
			wantSystemTrxData: &default_system.CSVSystemTrxData{
				TrxID:           "006630c83821fac6bea13b92b480feb2",
				TransactionTime: "2025-03-06 17:09:21",
				Type:            "DEBIT",
				Amount:          41000,
			},
			wantBankTrxData: &entitydefaultbank.CSVBankTrxData{
				DefaultUniqueIdentifier: "foo-5585fa85a971917b48ea2729bcf7d9fb",
				DefaultDate:             "2025-03-06",
				DefaultAmount:           -41000,
				DefaultBank:             "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Svc{
				comp: tt.fields.comp,
				repo: tt.fields.repo,
			}

			gotSystemTrxData, gotBankTrxData := s.parse(tt.args.data)
			if !reflect.DeepEqual(gotSystemTrxData, tt.wantSystemTrxData) {
				t.Errorf("parse() gotSystemTrxData = %v, want %v", gotSystemTrxData, tt.wantSystemTrxData)
			}

			if !reflect.DeepEqual(gotBankTrxData, tt.wantBankTrxData) {
				t.Errorf("parse() gotBankTrxData = %v, want %v", gotBankTrxData, tt.wantBankTrxData)
			}
		})
	}
}

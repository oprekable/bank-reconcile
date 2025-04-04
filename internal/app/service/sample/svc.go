package sample

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/repository/sample"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	entitybca "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bca/entity"
	entitybni "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/bni/entity"
	entitydefaultbank "github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks/default_bank/entity"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems/default_system"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/csvhelper"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/progressbarhelper"

	"github.com/aaronjan/hunch"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
	"go.chromium.org/luci/common/clock"
)

type Svc struct {
	comp *component.Components
	repo *repository.Repositories
}

var _ Service = (*Svc)(nil)

func NewSvc(
	comp *component.Components,
	repo *repository.Repositories,
) *Svc {
	return &Svc{
		comp: comp,
		repo: repo,
	}
}

func (s *Svc) deleteDirectorySystemTrxBankTrx(ctx context.Context, fs afero.Fs, isDeleteDirectory bool) (err error) {
	if !isDeleteDirectory {
		return
	}

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			e = csvhelper.DeleteDirectory(ctx, fs, s.comp.Config.Data.Reconciliation.SystemTRXPath)
			log.Err(ctx, "[sample.NewSvc] DeleteDirectory SystemTRXPath", e)
			return
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			e = csvhelper.DeleteDirectory(ctx, fs, s.comp.Config.Data.Reconciliation.BankTRXPath)
			log.Err(ctx, "[sample.NewSvc] DeleteDirectory BankTRXPath", e)
			return
		},
	)

	return
}

func (s *Svc) parse(data sample.TrxData) (systemTrxData systems.SystemTrxDataInterface, bankTrxData banks.BankTrxDataInterface) {
	if data.IsSystemTrx {
		systemTrxData = &default_system.CSVSystemTrxData{
			TrxID:           data.TrxID,
			TransactionTime: data.TransactionTime,
			Type:            data.Type,
			Amount:          data.Amount,
		}
	}

	if data.IsBankTrx || (!data.IsBankTrx && !data.IsSystemTrx) {
		bank := strings.ToLower(data.Bank)
		multiplier := float64(1)
		if data.Type == DEBIT {
			multiplier = float64(-1)
		}

		switch strings.ToUpper(bank) {
		case "BCA":
			{
				bankTrxData = &entitybca.CSVBankTrxData{
					BCAUniqueIdentifier: data.UniqueIdentifier,
					BCADate:             data.Date,
					BCAAmount:           data.Amount * multiplier,
					BCABank:             bank,
				}
			}
		case "BNI":
			{
				bankTrxData = &entitybni.CSVBankTrxData{
					BNIUniqueIdentifier: data.UniqueIdentifier,
					BNIDate:             data.Date,
					BNIAmount:           data.Amount * multiplier,
					BNIBank:             bank,
				}
			}
		default:
			{
				bankTrxData = &entitydefaultbank.CSVBankTrxData{
					DefaultUniqueIdentifier: data.UniqueIdentifier,
					DefaultDate:             data.Date,
					DefaultAmount:           data.Amount * multiplier,
					DefaultBank:             bank,
				}
			}
		}
	}

	return
}

func (s *Svc) GenerateSample(ctx context.Context, fs afero.Fs, bar *progressbar.ProgressBar, isDeleteDirectory bool) (returnSummary Summary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Sample Service").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())

	var trxData []sample.TrxData
	defer func() {
		_ = s.repo.RepoSample.Close()
		progressbarhelper.BarClear(bar)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return nil, s.deleteDirectorySystemTrxBankTrx(c, fs, isDeleteDirectory)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][1/5] Pre Process Generate Sample...")
			defer func() {
				log.Err(c, "[sample.NewSvc] RepoSample.Pre executed", e)
			}()

			e = s.repo.RepoSample.Pre(
				c,
				s.comp.Config.Data.Reconciliation.ListBank,
				s.comp.Config.Data.Reconciliation.FromDate,
				s.comp.Config.Data.Reconciliation.ToDate,
				s.comp.Config.Data.Reconciliation.TotalData,
				s.comp.Config.Data.Reconciliation.PercentageMatch,
			)

			return nil, e
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][2/5] Populate Trx Data...")

			trxData, e = s.repo.RepoSample.GetTrx(
				c,
			)

			log.Err(c, "[sample.NewSvc] RepoSample.GetTrx executed", e)
			return nil, e
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][3/5] Post Process Generate Sample...")
			if !s.comp.Config.Data.IsDebug {
				e = s.repo.RepoSample.Post(
					c,
				)
			}

			log.Err(c, "[sample.NewSvc] RepoSample.Post executed", e)

			return nil, e
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			progressbarhelper.BarDescribe(bar, "[cyan][4/5] Parse Sample Data...")
			systemTrxData := make([]systems.SystemTrxDataInterface, 0, len(trxData))
			bankTrxData := make(map[string][]banks.BankTrxDataInterface)

			lo.ForEach(trxData, func(data sample.TrxData, _ int) {
				systemTrx, bankTrx := s.parse(data)
				if systemTrx != nil {
					systemTrxData = append(systemTrxData, systemTrx)
				}

				if bankTrx != nil {
					bankTrxData[bankTrx.GetBank()] = append(bankTrxData[bankTrx.GetBank()], bankTrx)
				}
			})

			trxData = nil

			log.Msg(c, "[sample.NewSvc] populate systemTrxData & bankTrxData executed")
			progressbarhelper.BarDescribe(bar, "[cyan][5/5] Export Sample Data to CSV files...")

			lengthSystemTrxData := len(systemTrxData)
			fileNameSuffix := strconv.FormatInt(clock.Get(c).Now().Unix(), 10)
			returnSummary.FileSystemTrx = fmt.Sprintf("%s/%s.csv", s.comp.Config.Data.Reconciliation.SystemTRXPath, fileNameSuffix)
			returnSummary.TotalSystemTrx = int64(lengthSystemTrxData)

			sd := make([]*default_system.CSVSystemTrxData, 0, lengthSystemTrxData)

			lo.ForEach(systemTrxData, func(data systems.SystemTrxDataInterface, _ int) {
				sd = append(sd, data.(*default_system.CSVSystemTrxData))
			})

			systemTrxData = nil

			executor := make([]hunch.Executable, 0, len(bankTrxData)+1)
			executor = append(
				executor,
				func(ct context.Context) (interface{}, error) {
					er := csvhelper.StructToCSVFile(
						ct,
						fs,
						returnSummary.FileSystemTrx,
						sd,
						isDeleteDirectory,
					)

					log.Err(c, "[sample.NewSvc] save csv file "+returnSummary.FileSystemTrx+" executed", er)

					return nil, er
				},
			)

			returnSummary.TotalBankTrx = make(map[string]int64)
			returnSummary.FileBankTrx = make(map[string]string)

			for bankName, bankTrx := range bankTrxData {
				returnSummary.FileBankTrx[bankName] = fmt.Sprintf("%s/%s/%s_%s.csv", s.comp.Config.Data.Reconciliation.BankTRXPath, bankName, bankName, fileNameSuffix)
				totalBankTrx, exec := s.appendExecutor(
					fs,
					returnSummary.FileBankTrx[bankName],
					bankTrx,
					isDeleteDirectory,
				)

				if exec == nil || totalBankTrx == 0 {
					continue
				}

				returnSummary.TotalBankTrx[bankName] = totalBankTrx
				executor = append(executor, exec)
			}

			bankTrxData = nil

			return hunch.All(
				c,
				executor...,
			)
		},
	)

	return
}

func (s *Svc) appendExecutor(fs afero.Fs, filePath string, trxDataSlice []banks.BankTrxDataInterface, isDeleteDirectory bool) (totalData int64, executor hunch.Executable) {
	if len(trxDataSlice) == 0 {
		return 0, nil
	}

	formatText := "[sample.NewSvc] save csv file %s executed"

	switch trxDataSlice[0].(type) {
	case *entitybca.CSVBankTrxData:
		{
			bd := make([]*entitybca.CSVBankTrxData, 0, len(trxDataSlice))

			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*entitybca.CSVBankTrxData))
			})

			totalData = int64(len(bd))
			executor = func(ct context.Context) (interface{}, error) {
				er := csvhelper.StructToCSVFile(
					ct,
					fs,
					filePath,
					bd,
					isDeleteDirectory,
				)

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	case *entitybni.CSVBankTrxData:
		{
			bd := make([]*entitybni.CSVBankTrxData, 0, len(trxDataSlice))
			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*entitybni.CSVBankTrxData))
			})
			totalData = int64(len(bd))
			executor = func(ct context.Context) (interface{}, error) {
				er := csvhelper.StructToCSVFile(
					ct,
					fs,
					filePath,
					bd,
					isDeleteDirectory,
				)

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	default:
		{
			bd := make([]*entitydefaultbank.CSVBankTrxData, 0, len(trxDataSlice))
			lo.ForEach(trxDataSlice, func(data banks.BankTrxDataInterface, _ int) {
				bd = append(bd, data.(*entitydefaultbank.CSVBankTrxData))
			})
			totalData = int64(len(bd))
			executor = func(ct context.Context) (interface{}, error) {
				er := csvhelper.StructToCSVFile(
					ct,
					fs,
					filePath,
					bd,
					isDeleteDirectory,
				)

				log.Err(ct, fmt.Sprintf(formatText, filePath), er)
				return nil, er
			}
		}
	}

	return
}

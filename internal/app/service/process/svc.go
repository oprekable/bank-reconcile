package process

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/repository/process"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems/default_system"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/csvhelper"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/progressbarhelper"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/afero"
	"github.com/ulule/deepcopier"
	"go.chromium.org/luci/common/clock"
)

type Svc struct {
	comp                 *component.Components
	repo                 *repository.Repositories
	parserRegistry       *banks.ParserRegistry
	regexCompileBankName *regexp.Regexp
}

var _ ServiceGenerator = (*Svc)(nil)

func NewSvc(
	comp *component.Components,
	repo *repository.Repositories,
	parserRegistry *banks.ParserRegistry,
) *Svc {
	return &Svc{
		comp:                 comp,
		repo:                 repo,
		parserRegistry:       parserRegistry,
		regexCompileBankName: regexp.MustCompile(`.*[\\/]+([^\\/]+)[\\/][^\\/]+\.csv$`),
	}
}

func (s *Svc) parseSystemTrxFile(ctx context.Context, afs afero.Fs, filePath string) (returnData []*systems.SystemTrxData, err error) {
	var f afero.File
	f, err = afs.Open(filePath)
	if err != nil {
		log.Err(ctx, "[process.NewSvc] parseSystemTrxFile afs.Open - '"+filePath+"'", err)
		return
	}

	defer func() {
		if f != nil {
			_ = f.Close()
		}

		log.Err(ctx, "[process.NewSvc] parseSystemTrxFile parse - '"+filePath+"' executed", err)
	}()

	var systemParser *default_system.SystemParser
	if systemParser, err = default_system.NewSystemParser(
		&default_system.CSVSystemTrxData{},
		csv.NewReader(f),
		true,
	); err == nil {
		returnData, err = systemParser.ToSystemTrxData(ctx, filePath)
	}

	return
}

func (s *Svc) parseSystemTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*systems.SystemTrxData, err error) {
	var filePathSystemTrx []string
	defer func() {
		log.Err(ctx, "[process.NewSvc] parseSystemTrxFiles executed", err)
	}()

	cleanPath := filepath.Clean(s.comp.Config.Data.Reconciliation.SystemTRXPath)
	if err = afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".csv" {
			filePathSystemTrx = append(
				filePathSystemTrx,
				path,
			)
		}

		return nil
	}); err == nil {
		sliceMutex := sync.Mutex{}
		wg := sync.WaitGroup{}

		parallel.ForEach(filePathSystemTrx, func(item string, _ int) {
			wg.Add(1)
			defer wg.Done()
			data, _ := s.parseSystemTrxFile(ctx, afs, item)
			sliceMutex.Lock()
			returnData = append(returnData, data...)
			sliceMutex.Unlock()
		})

		wg.Wait()
	}

	return
}

func (s *Svc) importReconcileSystemDataToDB(ctx context.Context, data []*systems.SystemTrxData) (err error) {
	defSize := len(data) / s.comp.Config.Data.Reconciliation.NumberWorker
	numBigger := len(data) - defSize*s.comp.Config.Data.Reconciliation.NumberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < s.comp.Config.Data.Reconciliation.NumberWorker; i++ {
		if i == numBigger {
			size--
			if size == 0 {
				break
			}
		}

		err = s.repo.RepoProcess.ImportSystemTrx(
			ctx,
			data[idx:idx+size],
			idx,
			idx+size,
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) importReconcileMapToDB(ctx context.Context, min float64, max float64) (err error) {
	max = max + 1
	numberWorker := float64(s.comp.Config.Data.Reconciliation.NumberWorker * 2)
	defSize := max / numberWorker
	size := defSize + 1

	for i, idx := 0.0, min; i < numberWorker; i++ {
		err = s.repo.RepoProcess.GenerateReconciliationMap(
			ctx,
			idx,
			idx+size,
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) parseBankTrxFile(ctx context.Context, afs afero.Fs, item FilePathBankTrx) (returnData []*banks.BankTrxData, err error) {
	var bankParser banks.ReconcileBankData
	var f afero.File
	bank := strings.ToUpper(item.Bank)

	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			f, e = afs.Open(item.FilePath)
			return
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			// Use the injected registry to get the correct parser
			bankParser, e = s.parserRegistry.GetParser(bank, f, true)
			return
		},
	)

	log.Err(ctx, "[process.NewSvc] parseBankTrxFile parse ("+bank+") - '"+item.Bank+"' executed", err)

	if err != nil {
		return
	}

	returnData, err = bankParser.ToBankTrxData(ctx, item.FilePath)
	log.Err(ctx, "[process.NewSvc] parseBankTrxFile parse.ToBankTrxData ("+bank+") executed", err)

	return
}

func (s *Svc) parseBankTrxFiles(ctx context.Context, afs afero.Fs) (returnData []*banks.BankTrxData, err error) {
	var filePathBankTrx []FilePathBankTrx
	cleanPath := filepath.Clean(s.comp.Config.Data.Reconciliation.BankTRXPath)

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			// scan only csv file with first folder as bank name, bank should in the list of accepted bank name
			er := afero.Walk(afs, cleanPath, func(path string, info fs.FileInfo, err error) (e error) {
				match := s.regexCompileBankName.FindStringSubmatch(path)
				if len(match) <= 1 {
					return
				}

				if slices.Contains(s.comp.Config.Data.Reconciliation.ListBank, match[1]) {
					filePathBankTrx = append(
						filePathBankTrx,
						FilePathBankTrx{
							Bank:     match[1],
							FilePath: path,
						},
					)
				}

				return nil
			})

			return nil, er
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			sliceMutex := sync.Mutex{}
			wg := sync.WaitGroup{}

			parallel.ForEach(filePathBankTrx, func(item FilePathBankTrx, _ int) {
				wg.Add(1)
				defer wg.Done()
				data, _ := s.parseBankTrxFile(c, afs, item)
				sliceMutex.Lock()
				returnData = append(returnData, data...)
				sliceMutex.Unlock()
			})

			wg.Wait()
			return nil, nil
		},
	)

	log.Err(ctx, "[process.NewSvc] parseBankTrxFiles executed", err)

	return
}

func (s *Svc) importReconcileBankDataToDB(ctx context.Context, data []*banks.BankTrxData) (err error) {
	numberWorker := s.comp.Config.Data.Reconciliation.NumberWorker * 2
	defSize := len(data) / numberWorker
	numBigger := len(data) - defSize*numberWorker
	size := defSize + 1

	for i, idx := 0, 0; i < numberWorker; i++ {
		if i == numBigger {
			size--
			if size == 0 {
				break
			}
		}

		err = s.repo.RepoProcess.ImportBankTrx(
			ctx,
			data[idx:idx+size],
			idx,
			idx+size,
		)

		if err != nil {
			return
		}

		idx += size
	}

	return
}

func (s *Svc) parse(ctx context.Context, afs afero.Fs) (trxData parser.TrxData, err error) {
	isOK := func(t, minDate, maxDate time.Time) bool {
		return (t.Equal(minDate) || t.After(minDate)) && t.Before(maxDate)
	}

	isOKCheck := func(timeToCheck time.Time) bool {
		return isOK(
			timeToCheck,
			s.comp.Config.Data.Reconciliation.FromDate,
			s.comp.Config.Data.Reconciliation.ToDate.AddDate(0, 0, 1),
		)
	}

	setMaxAmount := func(currentAmount float64) {
		if trxData.MaxSystemAmount < currentAmount {
			trxData.MaxSystemAmount = currentAmount
		}
	}

	_, err = hunch.All(
		ctx,
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseSystemTrxFiles executed", e)
			}()

			var data []*systems.SystemTrxData
			if data, e = s.parseSystemTrxFiles(ct, afs); e == nil {
				trxData.SystemTrx = lo.Filter(data, func(item *systems.SystemTrxData, index int) bool {
					if !isOKCheck(item.TransactionTime) {
						return false
					}

					setMaxAmount(item.Amount)

					return true
				})
			}

			return
		},
		func(ct context.Context) (d interface{}, e error) {
			defer func() {
				log.Err(ct, "[process.NewSvc] GenerateReconciliation parseBankTrxFiles executed", e)
			}()

			var data []*banks.BankTrxData
			if data, e = s.parseBankTrxFiles(ct, afs); e == nil {
				trxData.BankTrx = lo.Filter(data, func(item *banks.BankTrxData, index int) bool {
					return isOKCheck(item.Date)
				})
			}

			return
		},
	)

	return
}

func (s *Svc) generateReconciliationSummaryAndFiles(ctx context.Context, fs afero.Fs, isDeleteDirectory bool) (returnData ReconciliationSummary, err error) {
	defer func() {
		log.Err(ctx, "[process.NewSvc] GenerateReconciliation RepoProcess.GetReconciliationSummary executed", err)
	}()

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			return s.repo.RepoProcess.GetReconciliationSummary(c)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			summary := i.(process.ReconciliationSummary)
			return nil, deepcopier.Copy(&summary).To(&returnData)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, s.generateReconciliationFiles(ctx, &returnData, fs, isDeleteDirectory)
		},
	)

	return
}

func (s *Svc) generateReconciliationFiles(ctx context.Context, reconciliationSummary *ReconciliationSummary, fs afero.Fs, isDeleteDirectory bool) (err error) {
	if reconciliationSummary == nil {
		return
	}

	fileNameSuffix := strconv.FormatInt(clock.Get(ctx).Now().Unix(), 10)
	logTemplate := "[process.NewSvc] save csv file %s executed"

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			fileName := fmt.Sprintf("%s/%s/%s/matched_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "system", "matched", fileNameSuffix)
			defer func() {
				log.Err(c, fmt.Sprintf(logTemplate, r), e)
			}()

			var d []process.MatchedTrx
			if d, e = s.repo.RepoProcess.GetMatchedTrx(ctx); e == nil && len(d) > 0 {
				return fileName, csvhelper.StructToCSVFile(
					c,
					fs,
					fileName,
					d,
					isDeleteDirectory,
				)
			}

			return "", e
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			reconciliationSummary.FileMatchedSystemTrx = i.(string)
			fileName := fmt.Sprintf("%s/%s/%s/not_matched_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "system", "not_matched", fileNameSuffix)
			defer func() {
				log.Err(c, fmt.Sprintf(logTemplate, r), e)
			}()

			var d []process.NotMatchedSystemTrx
			if d, e = s.repo.RepoProcess.GetNotMatchedSystemTrx(ctx); e == nil && len(d) > 0 {
				return fileName, csvhelper.StructToCSVFile(
					c,
					fs,
					fileName,
					d,
					isDeleteDirectory,
				)
			}

			return "", e
		},
		func(c context.Context, i interface{}) (_ interface{}, e error) {
			reconciliationSummary.FileMissingSystemTrx = i.(string)
			var d []process.NotMatchedBankTrx
			if d, e = s.repo.RepoProcess.GetNotMatchedBankTrx(ctx); e != nil || len(d) == 0 {
				return nil, e
			}

			bankTrxData := make(map[string][]process.NotMatchedBankTrx)
			lo.ForEach(d, func(data process.NotMatchedBankTrx, _ int) {
				data.Bank = strings.ToLower(data.Bank)
				bankTrxData[data.Bank] = append(bankTrxData[data.Bank], data)
			})

			reconciliationSummary.FileMissingBankTrx = make(map[string]string)
			bankNames := lo.Keys(bankTrxData)
			if isDeleteDirectory {
				if e = csvhelper.DeleteDirectory(
					c,
					fs,
					fmt.Sprintf("%s/%s/%s", s.comp.Config.Data.Reconciliation.ReportTRXPath, "bank", "not_matched"),
				); e != nil {
					return nil, e
				}
			}

			parallel.ForEach(bankNames, func(item string, _ int) {
				fileReportBankTrx := fmt.Sprintf("%s/%s/%s/%s_%s.csv", s.comp.Config.Data.Reconciliation.ReportTRXPath, "bank", "not_matched", item, fileNameSuffix)
				reconciliationSummary.FileMissingBankTrx[item] = fileReportBankTrx
				log.Err(
					c,
					fmt.Sprintf(logTemplate, fileReportBankTrx),
					csvhelper.StructToCSVFile(
						c,
						fs,
						fileReportBankTrx,
						bankTrxData[item],
						false,
					),
				)
			})

			return nil, e
		},
	)

	return
}

func (s *Svc) GenerateReconciliation(ctx context.Context, afs afero.Fs, bar *progressbar.ProgressBar) (returnData ReconciliationSummary, err error) {
	ctx = s.comp.Logger.GetLogger().With().Str("component", "Process ServiceGenerator").Ctx(ctx).Logger().WithContext(s.comp.Logger.GetCtx())

	defer func() {
		_ = s.repo.RepoProcess.Close()
		progressbarhelper.BarClear(bar)
	}()

	var trxData parser.TrxData

	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][1/7] Pre Process Generate Reconciliation...")

			e = s.repo.RepoProcess.Pre(
				c,
				s.comp.Config.Data.Reconciliation.ListBank,
				s.comp.Config.Data.Reconciliation.FromDate,
				s.comp.Config.Data.Reconciliation.ToDate,
			)

			log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Pre executed", e)
			return
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			progressbarhelper.BarDescribe(bar, "[cyan][2/7] Parse System/Bank Trx Files...")
			return s.parse(c, afs)
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][3/7] Import System Trx to DB...")

			trxData = i.(parser.TrxData)

			if len(trxData.SystemTrx) > 0 {
				e = s.importReconcileSystemDataToDB(c, trxData.SystemTrx)
			}

			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileSystemDataToDB executed", e)
			return
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][4/7] Import Bank Trx to DB...")

			if len(trxData.BankTrx) > 0 {
				e = s.importReconcileBankDataToDB(c, trxData.BankTrx)
			}

			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileBankDataToDB executed", e)

			return
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][5/7] Mapping Reconciliation Data...")

			if len(trxData.SystemTrx) > 0 {
				e = s.importReconcileMapToDB(c, trxData.MinSystemAmount, trxData.MaxSystemAmount)
			}

			log.Err(c, "[process.NewSvc] GenerateReconciliation importReconcileMapToDB executed", e)

			return
		},
		func(c context.Context, i interface{}) (d interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][6/7] Generate Reconciliation Report Files...")
			defer func() {
				log.Err(c, "[process.NewSvc] GenerateReconciliation generateReconciliationSummaryAndFiles executed", e)
			}()

			returnData, e = s.generateReconciliationSummaryAndFiles(c, afs, s.comp.Config.IsDeleteCurrentReportDirectory)
			return
		},
		func(c context.Context, i interface{}) (r interface{}, e error) {
			progressbarhelper.BarDescribe(bar, "[cyan][7/8] Post Process Generate Reconciliation...")
			if !s.comp.Config.Data.IsDebug {
				e = s.repo.RepoProcess.Post(
					c,
				)
				log.Err(c, "[process.NewSvc] GenerateReconciliation RepoProcess.Post executed", e)
			}

			return nil, e
		},
	)

	return
}

package process

import (
	"context"
	"fmt"
	"io"

	"github.com/aaronjan/hunch"
	"github.com/dustin/go-humanize"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/helper"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/memstats"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/tablewriterhelper"
)

const name = "process"

type Handler struct {
	comp   *component.Components
	svc    *service.Services
	repo   *repository.Repositories
	writer io.Writer
}

func NewHandler(writer io.Writer) *Handler {
	return &Handler{
		writer: writer,
	}
}

func (h *Handler) Name() string {
	return name
}

func (h *Handler) SetComponents(c *component.Components) {
	h.comp = c
}
func (h *Handler) SetServices(s *service.Services) {
	h.svc = s
}
func (h *Handler) SetRepositories(r *repository.Repositories) {
	h.repo = r
}

func (h *Handler) Exec() error {
	if h.comp == nil || h.svc == nil || h.repo == nil {
		return nil
	}

	var summary process.ReconciliationSummary
	bar := helper.InitProgressBar(h.writer)

	_, err := hunch.Waterfall(
		h.comp.Context,
		// Display config information
		func(c context.Context, _ interface{}) (interface{}, error) {
			formatText := "-%s --%s"
			args := helper.InitCommonArgs(
				h.comp.Config.Data,
				[][]string{
					{
						fmt.Sprintf(formatText, cmd.FlagReportTRXPathShort, cmd.FlagReportTRXPath),
						h.comp.Config.Data.Reconciliation.ReportTRXPath,
					},
				},
			)

			tableArgs := tablewriterhelper.InitTableWriter(h.writer)
			tableArgs.Header([]string{"Config", "Value"})
			_ = tableArgs.Bulk(args)
			_ = tableArgs.Render()

			return fmt.Fprintln(h.writer, "")
		},
		// Do reconcile process dan return summary data
		func(c context.Context, _ interface{}) (interface{}, error) {
			return h.svc.SvcProcess.GenerateReconciliation(h.comp.Context, h.comp.Fs.LocalStorageFs, bar)
		},
		// Display summary information
		func(c context.Context, i interface{}) (interface{}, error) {
			summary = i.(process.ReconciliationSummary)
			numberIntegerFormat := "#.###,"
			numberFloatFormat := "#.###,##"
			dataDesc := [][]string{
				{"Total number of transactions processed", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalProcessedSystemTrx))},
				{"Total number of matched transactions", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalMatchedSystemTrx))},
				{"Total number of not matched transactions", humanize.FormatInteger(numberIntegerFormat, int(summary.TotalNotMatchedSystemTrx))},
				{"Sum amount all transactions", humanize.FormatFloat(numberFloatFormat, summary.SumAmountProcessedSystemTrx)},
				{"Sum amount matched transactions", humanize.FormatFloat(numberFloatFormat, summary.SumAmountMatchedSystemTrx)},
				{"Total discrepancies", humanize.FormatFloat(numberFloatFormat, summary.SumAmountDiscrepanciesSystemTrx)},
			}

			_, _ = fmt.Fprintln(h.writer, "")
			tableDesc := tablewriterhelper.InitTableWriter(h.writer)
			tableDesc.Header([]string{"Description", "Value"})
			_ = tableDesc.Bulk(dataDesc)
			_ = tableDesc.Render()

			return fmt.Fprintln(h.writer, "")
		},
		// Display reconcile output files information
		func(c context.Context, i interface{}) (interface{}, error) {
			dataFilePath := [][]string{
				{"Matched system transaction data", summary.FileMatchedSystemTrx},
			}

			if summary.FileMissingSystemTrx != "" {
				dataFilePath = append(
					dataFilePath,
					[]string{"Missing system transaction data", summary.FileMissingSystemTrx},
				)
			}

			for bank, value := range summary.FileMissingBankTrx {
				dataFilePath = append(
					dataFilePath,
					[]string{
						fmt.Sprintf("Missing bank statement data - %s", bank),
						value,
					},
				)
			}

			_, _ = fmt.Fprintln(h.writer, "")
			tableFilePath := tablewriterhelper.InitTableWriter(h.writer)
			tableFilePath.Header([]string{"Description", "File Path"})
			_ = tableFilePath.Bulk(dataFilePath)
			_ = tableFilePath.Render()
			return fmt.Fprintln(h.writer, "")
		},
		// Display memory information
		func(c context.Context, i interface{}) (interface{}, error) {
			bar.Describe("[cyan]Done")
			memstats.PrintMemoryStats(h.writer)
			return nil, nil
		},
	)

	return err
}

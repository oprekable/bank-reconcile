package process

import (
	"context"
	"fmt"
	"io"

	"github.com/oprekable/bank-reconcile/cmd"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/helper"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/memstats"

	"github.com/dustin/go-humanize"

	"github.com/olekukonko/tablewriter"
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

func (h *Handler) Exec() error {
	if h.comp == nil || h.svc == nil || h.repo == nil {
		return nil
	}
	bar := helper.InitProgressBar(h.writer)
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

	tableArgs := tablewriter.NewWriter(h.writer)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()
	_, _ = fmt.Fprintln(h.writer, "")

	summary, err := h.svc.SvcProcess.GenerateReconciliation(context.Background(), h.comp.Fs.LocalStorageFs, bar)
	if err != nil {
		return err
	}

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
	tableDesc := tablewriter.NewWriter(h.writer)
	tableDesc.SetHeader([]string{"Description", "Value"})
	tableDesc.SetBorder(false)
	tableDesc.SetAlignment(tablewriter.ALIGN_LEFT)
	tableDesc.SetAutoWrapText(false)
	tableDesc.AppendBulk(dataDesc)
	tableDesc.Render()
	_, _ = fmt.Fprintln(h.writer, "")

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
	tableFilePath := tablewriter.NewWriter(h.writer)
	tableFilePath.SetHeader([]string{"Description", "File Path"})
	tableFilePath.SetBorder(false)
	tableFilePath.SetAlignment(tablewriter.ALIGN_LEFT)
	tableFilePath.SetAutoWrapText(false)
	tableFilePath.AppendBulk(dataFilePath)
	tableFilePath.Render()
	_, _ = fmt.Fprintln(h.writer, "")

	bar.Describe("[cyan]Done")
	memstats.PrintMemoryStats(h.writer)

	return nil
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

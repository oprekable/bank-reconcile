package process

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/helper"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/memstats"
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

	tableArgs := tablewriter.NewTable(
		h.writer,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: tw.NewSymbols(tw.StyleASCII),
				Settings: tw.Settings{
					Separators: tw.Separators{BetweenRows: tw.On},
					Lines:      tw.Lines{ShowFooterLine: tw.On},
				},
			},
		)),
		tablewriter.WithConfig(
			tablewriter.Config{
				Row: tw.CellConfig{
					Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				},
			},
		),
	)

	tableArgs.Header([]string{"Config", "Value"})
	_ = tableArgs.Bulk(args)
	_ = tableArgs.Render()
	_, _ = fmt.Fprintln(h.writer, "")

	summary, err := h.svc.SvcProcess.GenerateReconciliation(h.comp.Context, h.comp.Fs.LocalStorageFs, bar)
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

	tableDesc := tablewriter.NewTable(
		h.writer,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: tw.NewSymbols(tw.StyleASCII),
				Settings: tw.Settings{
					Separators: tw.Separators{BetweenRows: tw.On},
					Lines:      tw.Lines{ShowFooterLine: tw.On},
				},
			},
		)),
		tablewriter.WithConfig(
			tablewriter.Config{
				Row: tw.CellConfig{
					Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				},
			},
		),
	)

	tableDesc.Header([]string{"Description", "Value"})
	_ = tableDesc.Bulk(dataDesc)
	_ = tableDesc.Render()
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
	tableFilePath := tablewriter.NewTable(
		h.writer,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: tw.NewSymbols(tw.StyleASCII),
				Settings: tw.Settings{
					Separators: tw.Separators{BetweenRows: tw.On},
					Lines:      tw.Lines{ShowFooterLine: tw.On},
				},
			},
		)),
		tablewriter.WithConfig(
			tablewriter.Config{
				Row: tw.CellConfig{
					Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				},
			},
		),
	)

	tableFilePath.Header([]string{"Description", "File Path"})
	_ = tableFilePath.Bulk(dataFilePath)
	_ = tableFilePath.Render()
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

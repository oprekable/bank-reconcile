package sample

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/helper"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/memstats"

	"github.com/olekukonko/tablewriter"
)

const name = "sample"

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

func (h *Handler) Exec() (err error) {
	if h.comp == nil || h.svc == nil || h.repo == nil {
		return nil
	}
	bar := helper.InitProgressBar(h.writer)
	formatText := "-%s --%s"
	args := helper.InitCommonArgs(
		[][]string{
			{
				fmt.Sprintf(formatText, root.FlagTotalDataSampleToGenerateShort, root.FlagTotalDataSampleToGenerate),
				strconv.FormatInt(root.FlagTotalDataSampleToGenerateValue, 10),
			},
			{
				fmt.Sprintf(formatText, root.FlagPercentageMatchSampleToGenerateShort, root.FlagPercentageMatchSampleToGenerate),
				strconv.Itoa(root.FlagPercentageMatchSampleToGenerateValue),
			},
			{
				fmt.Sprintf(formatText, root.FlagIsDeleteCurrentSampleDirectoryShort, root.FlagIsDeleteCurrentSampleDirectory),
				strconv.FormatBool(root.FlagIsDeleteCurrentSampleDirectoryValue),
			},
		},
	)

	_, _ = fmt.Fprintln(h.writer, "")
	tableArgs := tablewriter.NewWriter(h.writer)
	tableArgs.SetHeader([]string{"Config", "Value"})
	tableArgs.SetBorder(false)
	tableArgs.SetAlignment(tablewriter.ALIGN_LEFT)
	tableArgs.AppendBulk(args)
	tableArgs.Render()

	summary, err := h.svc.SvcSample.GenerateSample(context.Background(), h.comp.Fs.LocalStorageFs, bar, h.comp.Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory)
	if err != nil {
		return err
	}

	data := [][]string{
		{"System Trx", "-", "Total Trx", strconv.FormatInt(summary.TotalSystemTrx, 10)}, //nolint:gofmt
		{"System Trx", "-", "File Path", summary.FileSystemTrx},
	}

	for k, v := range summary.FileBankTrx {
		data = append(
			data,
			[]string{"Bank Trx", k, "Total Trx", strconv.FormatInt(summary.TotalBankTrx[k], 10)},
			[]string{"Bank Trx", k, "File Path", v},
		)
	}

	_, _ = fmt.Fprintln(h.writer, "")
	table := tablewriter.NewWriter(h.writer)
	table.SetHeader([]string{"Type Trx", "Bank", "Title", ""})
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
	_, _ = fmt.Fprintln(h.writer, "")

	bar.Describe("[cyan]Done")
	memstats.PrintMemoryStats(h.writer)

	return
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

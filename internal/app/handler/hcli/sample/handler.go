package sample

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/oprekable/bank-reconcile/cmd"
	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/handler/hcli/helper"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
	"github.com/oprekable/bank-reconcile/internal/app/service/sample"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/memstats"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/tablewriterhelper"
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
		h.comp.Config.Data,
		[][]string{
			{
				fmt.Sprintf(formatText, cmd.FlagTotalDataSampleToGenerateShort, cmd.FlagTotalDataSampleToGenerate),
				strconv.FormatInt(h.comp.Config.Data.Reconciliation.TotalData, 10),
			},
			{
				fmt.Sprintf(formatText, cmd.FlagPercentageMatchSampleToGenerateShort, cmd.FlagPercentageMatchSampleToGenerate),
				strconv.Itoa(h.comp.Config.Data.Reconciliation.PercentageMatch),
			},
			{
				fmt.Sprintf(formatText, cmd.FlagIsDeleteCurrentSampleDirectoryShort, cmd.FlagIsDeleteCurrentSampleDirectory),
				strconv.FormatBool(h.comp.Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory),
			},
		},
	)

	_, _ = fmt.Fprintln(h.writer, "")
	tableArgs := tablewriterhelper.InitTableWriter(h.writer)
	tableArgs.Header([]string{"Config", "Value"})
	_ = tableArgs.Bulk(args)
	_ = tableArgs.Render()

	var summary sample.Summary
	summary, err = h.svc.SvcSample.GenerateSample(
		h.comp.Context,
		h.comp.Fs.LocalStorageFs,
		bar,
		h.comp.Config.Data.Reconciliation.IsDeleteCurrentSampleDirectory,
	)

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

	table := tablewriterhelper.InitTableWriter(h.writer)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Row.Formatting = tw.CellFormatting{MergeMode: tw.MergeHierarchical}
	})

	table.Header([]string{"Type Trx", "Bank", "Title", ""})
	_ = table.Bulk(data)
	_ = table.Render()
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

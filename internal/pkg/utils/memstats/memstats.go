package memstats

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func PrintMemoryStats(w io.Writer) {
	mem := MemStats()
	data := [][]string{
		{
			"Allocated", humanize.Bytes(mem.Alloc),
		},
		{
			"Total Allocated", humanize.Bytes(mem.TotalAlloc),
		},
		{
			"Memory Allocations", humanize.Bytes(mem.Mallocs),
		},
		{
			"Memory Frees", humanize.Bytes(mem.Frees),
		},
		{
			"Heap Allocated", humanize.Bytes(mem.HeapAlloc),
		},
		{
			"Heap System", humanize.Bytes(mem.HeapSys),
		},
		{
			"Heap In Use", humanize.Bytes(mem.HeapInuse),
		},
		{
			"Heap Idle", humanize.Bytes(mem.HeapIdle),
		},
		{
			"Heap OS Related", humanize.Bytes(mem.HeapReleased),
		},
		{
			"Heap Objects", humanize.Bytes(mem.HeapObjects),
		},
		{
			"Stack In Use", humanize.Bytes(mem.StackInuse),
		},
		{
			"Stack System", humanize.Bytes(mem.StackSys),
		},
		{
			"Stack Span In Use", humanize.Bytes(mem.MSpanInuse),
		},
		{
			"Stack Cache In Use", humanize.Bytes(mem.MCacheInuse),
		},
		{
			"Next GC cycle", humanizeNano(mem.NextGC),
		},
		{
			"Last GC cycle", humanize.Time(time.Unix(0, int64(mem.LastGC))), //nolint:gosec
		},
	}

	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "-------- Memory Dump --------")
	_, _ = fmt.Fprintln(w, "")

	tableDesc := tablewriter.NewTable(
		w,
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
	_ = tableDesc.Bulk(data)
	_ = tableDesc.Render()
	_, _ = fmt.Fprintln(w, "")
}

func MemStats() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}

func humanizeNano(n uint64) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix //nolint:gosec
}

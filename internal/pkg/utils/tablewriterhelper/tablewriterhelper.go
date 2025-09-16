package tablewriterhelper

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func InitTableWriter(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewTable(
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

	return table
}

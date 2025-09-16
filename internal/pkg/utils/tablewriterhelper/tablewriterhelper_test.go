package tablewriterhelper

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestInitTableWriter(t *testing.T) {
	tests := []struct {
		name  string
		data  [][]string
		wantW []string
	}{
		{
			name: "Ok",
			data: [][]string{
				{"foo", "fofofo"},
				{"bar", "barbarbar"},
			},
			wantW: []string{
				"DESCRIPTION",
				"fofofo",
				"barbarbar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			table := InitTableWriter(w)
			table.Header([]string{"Description", "Value"})
			_ = table.Bulk(tt.data)
			_ = table.Render()
			_, _ = fmt.Fprintln(w, "")

			got := w.String()
			for _, want := range tt.wantW {
				if !strings.Contains(got, want) {
					t.Errorf("InitTableWriter() = %v, want %v", got, tt.wantW)
					break
				}
			}
			w.Reset()
		})
	}
}

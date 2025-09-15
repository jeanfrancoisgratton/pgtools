// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/15 08:08
// Original filename: src/conf/renderer.go

package conf

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Render prints a compact table, similar styling to ListEnvironments().
// Also truncates the Description column to 40 characters with "...".
func Render(rows []Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Setting", "Unit", "Source", "Category", "Description"})
	for _, r := range rows {
		cat := ellipsize(r.Category, 30)
		desc := ellipsize(r.ShortDesc, 40)
		t.AppendRow(table.Row{r.Name, r.Setting, r.Unit, r.Source, cat, desc})
	}
	t.SortBy([]table.SortBy{{Name: "Name", Mode: table.Asc}})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()
}

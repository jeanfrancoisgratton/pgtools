// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:38
// Original filename: src/roles/helpers.go

package roles

import (
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// quoteIdent applies a minimal PostgreSQL identifier quoting.
func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func buildHeader(includeMembers, verbose bool) table.Row {
	cols := []any{"Role"}
	if verbose {
		cols = append(cols, "Login", "Superuser", "CreateDB", "CreateRole", "Inherit", "Replication", "BypassRLS")
	} else {
		cols = append(cols, "Login", "Superuser", "CreateDB", "CreateRole")
	}
	if includeMembers {
		cols = append(cols, "Members")
	}
	return table.Row(cols)
}

func buildRow(r roleRow, includeMembers, verbose bool, members map[string][]string) table.Row {
	cols := []any{r.Name}
	if verbose {
		cols = append(cols, yesNo(r.Login), yesNo(r.Superuser), yesNo(r.CreateDB), yesNo(r.CreateRole),
			yesNo(r.Inherit), yesNo(r.Replication), yesNo(r.BypassRLS))
	} else {
		cols = append(cols, yesNo(r.Login), yesNo(r.Superuser), yesNo(r.CreateDB), yesNo(r.CreateRole))
	}
	if includeMembers {
		ms := members[r.Name]
		sort.Strings(ms)
		cols = append(cols, strings.Join(ms, ", "))
	}
	return table.Row(cols)
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

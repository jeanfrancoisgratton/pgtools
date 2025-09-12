// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 03:58
// Original filename: src/roles/sharedFunctions.go

package shared

import "strings"

// QuoteIdent applies a minimal PostgreSQL identifier quoting.
func QuoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

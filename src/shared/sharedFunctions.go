// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 03:58
// Original filename: src/roles/sharedFunctions.go

package shared

import (
	"fmt"
	"strings"
)

// QuoteIdent quotes a single identifier with double quotes, escaping any internal quotes.
func QuoteIdent(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}

// QuoteQualifiedIdent quotes a schema and table as "schema"."table".
func QuoteQualifiedIdent(schema, table string) string {
	return QuoteIdent(schema) + "." + QuoteIdent(table)
}

// QuoteIdents applies QuoteIdent to a slice of identifiers.
func QuoteIdents(in []string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = QuoteIdent(s)
	}
	return out
}

// HumanizeBytes renders a byte count using 1024-based units with a compact format.
// Examples: 0B, 999B, 1.0KiB, 12.3MiB, 1.0GiB
func HumanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit && exp < 5; n /= unit {
		div *= unit
		exp++
	}
	// Units: KiB, MiB, GiB, TiB, PiB
	pre := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}[exp]
	value := float64(b) / float64(div)
	return fmt.Sprintf("%.3f%s", value, pre)
}

// HumanizeBytesMBGB formats a byte count into MB or GB with 1 decimal place.
// It deliberately avoids KB/TB to keep the output compact and aligned.
func HumanizeBytesMBGB(b int64) string {
	const (
		MB = 1 << 20
		GB = 1 << 30
	)

	if b < GB {
		mb := float64(b) / float64(MB)
		return fmt.Sprintf("%.1f MB", mb)
	}
	gb := float64(b) / float64(GB)
	return fmt.Sprintf("%.1f GB", gb)
}

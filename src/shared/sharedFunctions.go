// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 03:58
// Original filename: src/roles/sharedFunctions.go

package shared

import (
	"fmt"
	"strings"
)

// QuoteIdent applies a minimal PostgreSQL identifier quoting.
func QuoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
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

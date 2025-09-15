// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 19:19
// Original filename: src/show/tableHelpers.go

package show

import (
	"fmt"
	"strings"
)

// ExcludePattern represents a parsed entry from -X/--exclude.
// Supported forms (case-insensitive, folded to lower):
//
//	"table"
//	"db.table"
//
// (optional future) "db.schema.table"
type ExcludePattern struct {
	DB     string // optional
	Schema string // optional (reserved for future)
	Table  string // required
}

// BuildTableExcluder returns a predicate usable as:
//
//	if excluder(db, schema, table) { /* skip */ }
func BuildTableExcluder(raw []string) (func(db, schema, table string) bool, error) {
	patterns := make([]ExcludePattern, 0, len(raw))

	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		p := ExcludePattern{}
		parts := strings.Split(s, ".")
		switch len(parts) {
		case 1:
			p.Table = strings.ToLower(parts[0])
		case 2:
			p.DB = strings.ToLower(parts[0])
			p.Table = strings.ToLower(parts[1])
		case 3:
			// Uncomment to support db.schema.table explicitly
			// p.DB    = strings.ToLower(parts[0])
			// p.Schema= strings.ToLower(parts[1])
			// p.Table = strings.ToLower(parts[2])
			return nil, fmt.Errorf("unsupported exclude pattern %q (use 'table' or 'db.table')", s)
		default:
			return nil, fmt.Errorf("invalid exclude pattern %q", s)
		}
		if p.Table == "" {
			return nil, fmt.Errorf("invalid exclude pattern %q: empty table", s)
		}
		patterns = append(patterns, p)
	}

	// Return predicate
	return func(db, schema, table string) bool {
		db = strings.ToLower(db)
		schema = strings.ToLower(schema)
		table = strings.ToLower(table)

		for _, p := range patterns {
			// Match DB if set
			if p.DB != "" && p.DB != db {
				continue
			}
			// (Schema-aware matching could go here later)
			if p.Table == table {
				return true
			}
		}
		return false
	}, nil
}

// InSetCI checks membership case-insensitively.
func InSetCI(needle string, hay []string) bool {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for _, h := range hay {
		if strings.ToLower(strings.TrimSpace(h)) == needle {
			return true
		}
	}
	return false
}

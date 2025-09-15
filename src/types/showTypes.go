// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 19:18
// Original filename: src/types/showTypes.go

package types

var SortBySize bool
var DB2Show = ""

type DbRow struct {
	Name      string
	SizeBytes int64
}

type SchemaRow struct {
	DB        string
	Schema    string
	Owner     string
	Tables    int64
	Views     int64
	TotalSize string // already pretty-printed (e.g., "17 MB")
}

// CLI-bound variables (with reasonable defaults).
var ExcludedDBs = []string{"template0", "template1"}
var ExcludedTables = []string{}

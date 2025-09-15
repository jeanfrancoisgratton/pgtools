// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 19:18
// Original filename: src/types/types_show.go

package show

import "time"

var SortBySize bool
var DefaultDB = "postgres"
var PagedOutput = false

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

type TableRow struct {
	DB      string
	Schema  string
	Table   string
	Owner   string
	RowsEst int64
	TotalB  int64
	TableB  int64
	IndexB  int64
	ToastB  int64
	HasPK   bool
}

// DbStatsRow holds the handful of counters we display.

type DbStatsRow struct {
	DB        string
	Txns      int64 // xact_commit + xact_rollback
	Commits   int64 // xact_commit
	Deadlocks int64 // deadlocks
	StatsAge  time.Duration
}

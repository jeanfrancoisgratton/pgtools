// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/15 05:13
// Original filename: src/show/showStats.go
// Purpose: "pgtools show stats" — summarize pg_stat_database counters.

package show

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowStats queries pg_stat_database and prints a compact summary.
// We deliberately omit template0/1. Sorted by Txns desc.
func ShowStats(ctx context.Context, pool *pgxpool.Pool) *ce.CustomError {
	const q = `
SELECT
  datname,
  xact_commit,
  xact_rollback,
  deadlocks,
  stats_since
FROM pg_stat_database
WHERE datname NOT IN ('template0','template1')
ORDER BY datname;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return &ce.CustomError{Code: 851, Title: "Query error", Message: err.Error()}
	}
	defer rows.Close()

	now := time.Now()
	var out []DbStatsRow

	for rows.Next() {
		var (
			db         string
			xCommit    int64
			xRollback  int64
			deadlocks  int64
			statsSince *time.Time
		)
		if scanErr := rows.Scan(&db, &xCommit, &xRollback, &deadlocks, &statsSince); scanErr != nil {
			return &ce.CustomError{Code: 852, Title: "Row scan error", Message: scanErr.Error()}
		}
		var age time.Duration
		if statsSince != nil && !statsSince.IsZero() {
			age = now.Sub(*statsSince)
			if age < 0 {
				age = 0
			}
		}
		out = append(out, DbStatsRow{
			DB:        db,
			Txns:      xCommit + xRollback,
			Commits:   xCommit,
			Deadlocks: deadlocks,
			StatsAge:  age,
		})
	}
	if err := rows.Err(); err != nil {
		return &ce.CustomError{Code: 853, Title: "Row iteration error", Message: err.Error()}
	}

	// Sort by Txns descending to bubble up the busiest DBs.
	sort.Slice(out, func(i, j int) bool { return out[i].Txns > out[j].Txns })

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{"Database", "Txns", "Commits", "Deadlocks", "Stats Age"})

	for _, r := range out {
		tw.AppendRow(table.Row{
			r.DB,
			r.Txns,
			r.Commits,
			r.Deadlocks,
			humanizeDurationCompact(r.StatsAge),
		})
	}

	tw.SetStyle(table.StyleBold)
	tw.Style().Format.Header = text.FormatDefault
	tw.Style().Color.Header = text.Colors{text.Bold}
	tw.Render()
	return nil
}

// humanizeDurationCompact prints durations like 3d4h12m or 45m or 8s.
// It keeps output compact for table alignment.
func humanizeDurationCompact(d time.Duration) string {
	if d <= 0 {
		return "0s"
	}
	// Round down to seconds to avoid noisy ms.
	sec := int64(d.Seconds())

	const (
		day  = int64(24 * 3600)
		hour = int64(3600)
		min  = int64(60)
	)

	days := sec / day
	sec %= day
	hrs := sec / hour
	sec %= hour
	mins := sec / min
	sec %= min

	str := ""
	if days > 0 {
		str += intToStr(days) + "d"
	}
	if hrs > 0 {
		str += intToStr(hrs) + "h"
	}
	if mins > 0 && days == 0 { // keep it compact: if days>0, hours are enough
		str += intToStr(mins) + "m"
	}
	if sec > 0 && days == 0 && hrs == 0 { // show seconds only for short ages
		str += intToStr(sec) + "s"
	}
	if str == "" {
		return "0s"
	}
	return str
}

func intToStr(v int64) string {
	// small helper to avoid fmt for tiny hot path
	// (consistency with other helpers that avoid excessive fmt usage)
	return string([]byte(int64ToBytes(v)))
}

func int64ToBytes(n int64) []byte {
	// minimal int64 -> bytes (base10) to avoid fmt in tight loops
	if n == 0 {
		return []byte{'0'}
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return buf[i:]
}

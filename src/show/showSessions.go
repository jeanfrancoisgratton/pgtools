// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 11:19
// Original filename: src/show/showSessions.go

package show

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListSessions prints current sessions (pg_stat_activity), including client IP.
func ShowSessions(ctx context.Context, pool *pgxpool.Pool) error {
	const q = `
SELECT
  pid,
  usename,
  datname,
  application_name,
  client_addr,
  client_port,
  backend_start,
  state,
  wait_event_type,
  wait_event,
  LEFT(COALESCE(query, ''), 120) AS query_snippet
FROM pg_stat_activity
ORDER BY backend_start;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return fmt.Errorf("query sessions: %w", err)
	}
	defer rows.Close()

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{
		"PID", "USER", "DB", "APP", "CLIENT_ADDR", "CLIENT_PORT", "STARTED", "STATE", "WAIT", "EVENT", "QUERY(120)",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "PID", Align: text.AlignRight},
		{Name: "CLIENT_PORT", Align: text.AlignRight},
		{Name: "STARTED", Align: text.AlignRight},
	})

	for rows.Next() {
		var (
			pid          int32
			user         string
			db           string
			app          string
			clientAddr   *string
			clientPort   *int32
			backendStart time.Time
			state        *string
			waitType     *string
			waitEvent    *string
			querySnippet string
		)
		if err := rows.Scan(&pid, &user, &db, &app, &clientAddr, &clientPort, &backendStart, &state, &waitType, &waitEvent, &querySnippet); err != nil {
			return fmt.Errorf("scan session: %w", err)
		}

		addr := ""
		if clientAddr != nil {
			addr = *clientAddr
		}
		port := ""
		if clientPort != nil {
			port = fmt.Sprintf("%d", *clientPort)
		}
		st := ""
		if state != nil {
			st = *state
		}
		wt := ""
		if waitType != nil {
			wt = *waitType
		}
		we := ""
		if waitEvent != nil {
			we = *waitEvent
		}

		tw.AppendRow(table.Row{
			pid,
			user,
			db,
			app,
			addr,
			port,
			backendStart.Format(time.RFC3339),
			st,
			wt,
			we,
			querySnippet,
		})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows (sessions): %w", err)
	}

	tw.Render()
	return nil
}

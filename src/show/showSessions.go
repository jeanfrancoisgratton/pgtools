// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 11:19
// Original filename: src/show/showSessions.go
// Updated: 2025/09/14  (NULL-safe scanning; Option B)

package show

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowSessions prints current sessions from pg_stat_activity.
// NOTE: This version is identical in behavior but uses sql.NullString / sql.NullInt32
// to avoid crashes when scanning NULLs from background/system backends.
func ShowSessions(ctx context.Context, pool *pgxpool.Pool) *ce.CustomError {
	const q = `SELECT pid,usename,datname,application_name,client_addr::text AS client_addr,client_port,backend_start,
       state,wait_event_type,wait_event FROM pg_stat_activity ORDER BY backend_start ASC;`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return &ce.CustomError{Code: 801, Title: "Query error", Message: err.Error()}
	}
	defer rows.Close()

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{
		"PID", "User", "DB", "App", "Client", "Port", "Started",
		"State", "WaitType", "WaitEvent",
	})

	for rows.Next() {
		var (
			pid          int32
			usename      sql.NullString
			datname      sql.NullString
			app          sql.NullString
			addr         sql.NullString
			port         sql.NullInt32
			backendStart time.Time
			state        sql.NullString
			waitType     sql.NullString
			waitEvent    sql.NullString
			//			query        sql.NullString
		)
		if err := rows.Scan(&pid, &usename, &datname, &app, &addr, &port, &backendStart, &state, &waitType, &waitEvent); err != nil {
			return &ce.CustomError{Code: 802, Title: "Error scanning sessions", Message: err.Error()}
		}

		userStr := ""
		if usename.Valid {
			userStr = usename.String
		}
		dbStr := ""
		if datname.Valid {
			dbStr = datname.String
		}
		appStr := ""
		if app.Valid {
			appStr = app.String
		}
		addrStr := ""
		if addr.Valid {
			addrStr = addr.String
		}
		stateStr := ""
		if state.Valid {
			stateStr = state.String
		}
		wtStr := ""
		if waitType.Valid {
			wtStr = waitType.String
		}
		weStr := ""
		if waitEvent.Valid {
			weStr = waitEvent.String
		}

		var portVal any = ""
		if port.Valid {
			portVal = port.Int32
		}

		tw.AppendRow(table.Row{
			pid, userStr, dbStr, appStr, addrStr, portVal,
			backendStart.Format(time.RFC3339),
			stateStr, wtStr, weStr,
		})
	}
	if err := rows.Err(); err != nil {
		return &ce.CustomError{Code: 803, Title: "Row scanning error", Message: err.Error()}
	}

	tw.SetStyle(table.StyleBold)
	tw.Style().Format.Header = text.FormatDefault
	tw.Style().Color.Header = text.Colors{text.Bold}
	tw.Render()
	return nil
}

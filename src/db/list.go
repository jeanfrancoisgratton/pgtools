// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/09 02:55
// Original filename: src/db/list.go

package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"pgtool/logging"
	"pgtool/types"
)

// List all databases on the server
func ListDatabases(cfg *types.DBConfig) ([]string, *ce.CustomError) {
	conn, cerr := Connect(cfg, "postgres")
	if cerr != nil {
		return nil, cerr
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(),
		"SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname")
	if err != nil {
		return nil, &ce.CustomError{Title: "Query failed", Message: err.Error()}
	}
	defer rows.Close()

	var dbs []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, &ce.CustomError{Title: "scan failed", Message: err.Error()}
		}
		dbs = append(dbs, name)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Title: "row error", Message: err.Error()}
	}

	if types.Quiet {
		return dbs, nil
	}

	// Now that we have the list, let's print it
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Database name"})

	for _, db := range dbs {
		t.AppendRow([]interface{}{hf.Green(db)})
	}
	t.SortBy([]table.SortBy{
		{Name: "Database name", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return dbs, nil
}

func ListTables(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: ListTables")
	rows, err := conn.Query(context.Background(), `
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public'
	`)
	if err != nil {
		return nil, &ce.CustomError{Code: 103, Title: "Unable to list tables", Message: err.Error()}
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, &ce.CustomError{Code: 104, Title: "Scan failed", Message: err.Error()}
		}
		tables = append(tables, t)
	}
	return tables, nil
}

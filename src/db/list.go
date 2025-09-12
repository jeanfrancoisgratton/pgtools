// db/list.go
// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Rewritten: 2025/09/12
// File: db/list.go
//
// Purpose:
//
//	ListDatabases() now accepts a sort-by-size switch. It fetches database
//	sizes via pg_database_size(datname), sorts either by name (default) or
//	by raw size bytes (when requested), prints a table (unless types.Quiet),
//	and returns the database names.
//
// Notes:
//   - Sorting by size uses raw bytes (accurate across MB/GB labels).
//   - Name sort is ascending; size sort is descending (largest first).
//   - Keeps output style consistent with the rest of the tool.
package db

import (
	"context"
	"os"
	"pgtools/shared"
	"sort"

	"pgtools/logging"
	"pgtools/types"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// dbRow holds a single database row from the list query.
type dbRow struct {
	Name      string
	SizeBytes int64
}

// ListDatabases connects to the server, retrieves non-template databases with
// their sizes, sorts by name (default) or size (if sortBySize == true), prints
// a table (unless types.Quiet), and returns the list of database names.
func ListDatabases(cfg *types.DBConfig, sortBySize bool) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: ListDatabases")
	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, qerr := conn.Query(
		context.Background(),
		`SELECT datname, pg_database_size(datname) AS size_bytes
		 FROM pg_database
		 WHERE datistemplate = false
		 ORDER BY datname`,
	)
	if qerr != nil {
		return nil, &ce.CustomError{Code: 101, Title: "Unable to list databases", Message: qerr.Error()}
	}
	defer rows.Close()

	var (
		dbs    []string
		result []dbRow
	)
	for rows.Next() {
		var name string
		var size int64
		if scanErr := rows.Scan(&name, &size); scanErr != nil {
			return nil, &ce.CustomError{Code: 102, Title: "Scan failed", Message: scanErr.Error()}
		}
		dbs = append(dbs, name)
		result = append(result, dbRow{Name: name, SizeBytes: size})
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, &ce.CustomError{Code: 103, Title: "List iteration failed", Message: rowsErr.Error()}
	}

	// Sort according to flag:
	if sortBySize {
		// Largest first
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].SizeBytes > result[j].SizeBytes
		})
	} else {
		// Name ascending
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
	}

	// Render table unless quiet.
	if !types.Quiet {
		renderDBTable(result)
	}

	// Return names in the same order as rendered:
	names := make([]string, len(result))
	for i, r := range result {
		names[i] = r.Name
	}
	return names, nil
}

// renderDBTable prints the database table using go-pretty with a Size column.
func renderDBTable(items []dbRow) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Database name", "Size"})

	for _, r := range items {
		t.AppendRow([]interface{}{
			hf.Green(r.Name),
			shared.HumanizeBytesMBGB(r.SizeBytes),
		})
	}

	// We already sorted prior to rendering; keep table style only.
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()
}

// ListTables returns all table names from the connected database.
// (Kept here because the original file also exposed ListTables.)
func ListTables(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: ListTables")
	rows, err := conn.Query(
		context.Background(),
		`SELECT table_name
		   FROM information_schema.tables
		  WHERE table_schema='public'
		  ORDER BY table_name`,
	)
	if err != nil {
		return nil, &ce.CustomError{Code: 104, Title: "Unable to list tables", Message: err.Error()}
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tname string
		if err := rows.Scan(&tname); err != nil {
			return nil, &ce.CustomError{Code: 105, Title: "Scan failed", Message: err.Error()}
		}
		tables = append(tables, tname)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 106, Title: "List iteration failed", Message: rows.Err().Error()}
	}
	return tables, nil
}

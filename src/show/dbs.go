// show/showDBs.go
// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Rewritten: 2025/09/12

package show

import (
	"context"
	"os"
	"sort"

	"pgtools/logging"
	"pgtools/shared"
	"pgtools/types"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DbRow holds one database entry (name + size in bytes).

// ShowDatabases connects to the server, retrieves non-template databases with
// their sizes, optionally sorts by size (largest first) when sortBySize is true,
// prints a table unless types.Quiet is set, and returns the show of database names.
func ShowDatabases(cfg *types.DBConfig, sortBySize bool) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: show.ShowDatabases")

	// Connect to the maintenance DB "postgres".
	dsn := shared.BuildDSN(cfg, "postgres")
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, &ce.CustomError{Title: "Connection failure", Message: err.Error(), Code: 200}
	}
	defer conn.Close(context.Background())

	rows, qerr := conn.Query(context.Background(), `
		SELECT datname, pg_database_size(datname) AS size_bytes
		FROM pg_database
		WHERE datistemplate = false
	`)
	if qerr != nil {
		return nil, &ce.CustomError{Title: "Unable to show databases", Message: qerr.Error(), Code: 201}
	}
	defer rows.Close()

	var data []DbRow
	for rows.Next() {
		var r DbRow
		if err := rows.Scan(&r.Name, &r.SizeBytes); err != nil {
			return nil, &ce.CustomError{Title: "Scan failed", Message: err.Error(), Code: 202}
		}
		data = append(data, r)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Title: "List iteration failed", Message: rows.Err().Error(), Code: 203}
	}

	// Sort
	if sortBySize {
		sort.Slice(data, func(i, j int) bool {
			return data[i].SizeBytes > data[j].SizeBytes // descending by size
		})
	} else {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name < data[j].Name // ascending by name
		})
	}

	// Print unless quiet
	if !types.Quiet {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Database", "Size"})
		t.SetStyle(table.StyleBold)
		t.Style().Format.Header = text.FormatDefault
		t.Style().Color.Header = text.Colors{text.Bold}

		for _, r := range data {
			t.AppendRow(table.Row{r.Name, shared.HumanizeBytes(r.SizeBytes)}) // uses helperFunctions HumanizeBytes
		}
		t.Render()
	}

	// Return names
	names := make([]string, 0, len(data))
	for _, r := range data {
		names = append(names, r.Name)
	}
	return names, nil
}

// ListTables returns a show of "schema.table" names for the connected DB.
// NOTE: signature matches usage in db/dump_data.go: ListTables(conn *pgx.Conn)
func ListTables(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: show.ListTables")

	rows, err := conn.Query(context.Background(), `
		SELECT schemaname, tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname NOT IN ('pg_catalog','information_schema')
		ORDER BY schemaname, tablename
	`)
	if err != nil {
		return nil, &ce.CustomError{Title: "Unable to show tables", Message: err.Error(), Code: 104}
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			return nil, &ce.CustomError{Title: "Scan failed", Message: err.Error(), Code: 105}
		}
		tables = append(tables, schema+"."+name)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Title: "List iteration failed", Message: rows.Err().Error(), Code: 106}
	}
	return tables, nil
}

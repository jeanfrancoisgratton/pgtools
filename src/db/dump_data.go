// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Rewritten: 2025/09/12
// Original filename: src/db/dump_data.go

package db

import (
	"context"
	"fmt"
	"io"
	"pgtools/logging"
	"pgtools/shared"
	"pgtools/types"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// writeDatabaseSQL connects to dbName, enumerates user tables, and writes INSERT statements
// for their contents to the provided writer.
// NOTE: This version intentionally avoids importing pgtools/show to break the package cycle.
func writeDatabaseSQL(cfg *types.DBConfig, dbName string, writer io.Writer) *ce.CustomError {
	logging.Debugf("Entering writeDatabaseSQL for %s", dbName)
	logging.Infof("Processing database: %s", dbName)

	// Connect to target DB
	conn, err := Connect(cfg, dbName)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// Optional header
	fmt.Fprintf(writer, "--\n-- Database: %s\n-- Generated at: %s\n--\n\n", dbName, time.Now().Format(time.RFC3339))
	fmt.Fprintln(writer, "BEGIN;")

	// List user tables
	tables, cerr := getTableNames(conn)
	if cerr != nil {
		return cerr
	}

	// Dump rows as INSERTs
	for _, fq := range tables {
		// fq is "schema.table"
		parts := strings.SplitN(fq, ".", 2)
		if len(parts) != 2 {
			return &ce.CustomError{Code: 204, Title: "Invalid table name", Message: "expected schema.table"}
		}
		schema, table := parts[0], parts[1]
		full := shared.QuoteQualifiedIdent(schema, table)

		// Fetch columns for this table
		cols, cerr := getColumnNames(conn, schema, table)
		if cerr != nil {
			return cerr
		}
		if len(cols) == 0 {
			// No columns? Skip.
			continue
		}

		// Construct SELECT
		selectSQL := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(shared.QuoteIdents(cols), ", "), full)
		ctx := context.Background()
		rows, qerr := conn.Query(ctx, selectSQL)
		if qerr != nil {
			return &ce.CustomError{Code: 205, Title: "Query failed", Message: qerr.Error()}
		}

		// Prepare INSERT prefix
		insertPrefix := fmt.Sprintf("INSERT INTO %s (%s) VALUES", full, strings.Join(shared.QuoteIdents(cols), ", "))

		// Stream rows → INSERT statements
		for rows.Next() {
			values, vErr := rows.Values()
			if vErr != nil {
				rows.Close()
				return &ce.CustomError{Code: 206, Title: "Row fetch failed", Message: vErr.Error()}
			}

			fmt.Fprint(writer, insertPrefix)
			fmt.Fprint(writer, " (")

			for i, val := range values {
				if i > 0 {
					fmt.Fprint(writer, ", ")
				}
				if val == nil {
					fmt.Fprint(writer, "NULL")
					continue
				}
				// Basic literal formatting + single-quote escaping
				escaped := strings.ReplaceAll(fmt.Sprint(val), "'", "''")
				fmt.Fprintf(writer, "'%s'", escaped)
			}

			fmt.Fprintln(writer, ");")
		}
		if rows.Err() != nil {
			rows.Close()
			return &ce.CustomError{Code: 207, Title: "Row iteration failed", Message: rows.Err().Error()}
		}
		rows.Close()
	}

	fmt.Fprintln(writer, "COMMIT;")
	return nil
}

// getTableNames lists user tables as "schema.table" in the connected DB.
func getTableNames(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: getTableNames")

	rows, err := conn.Query(context.Background(), `
		SELECT schemaname, tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname NOT IN ('pg_catalog','information_schema')
		ORDER BY schemaname, tablename`)
	if err != nil {
		return nil, &ce.CustomError{Code: 201, Title: "Unable to show tables", Message: err.Error()}
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var s, n string
		if err := rows.Scan(&s, &n); err != nil {
			return nil, &ce.CustomError{Code: 202, Title: "Scan error", Message: err.Error()}
		}
		out = append(out, s+"."+n)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 203, Title: "List iteration failed", Message: rows.Err().Error()}
	}
	return out, nil
}

// getColumnNames retrieves column names for a given table.
func getColumnNames(conn *pgx.Conn, schema, table string) ([]string, *ce.CustomError) {
	q := `
		SELECT a.attname
		FROM pg_attribute a
		JOIN pg_class c ON a.attrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = $1
		  AND c.relname = $2
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		ORDER BY a.attnum`
	rows, err := conn.Query(context.Background(), q, schema, table)
	if err != nil {
		return nil, &ce.CustomError{Code: 208, Title: "Column query failed", Message: err.Error()}
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, &ce.CustomError{Code: 209, Title: "Column scan failed", Message: err.Error()}
		}
		cols = append(cols, name)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 210, Title: "Column iteration failed", Message: rows.Err().Error()}
	}
	return cols, nil
}

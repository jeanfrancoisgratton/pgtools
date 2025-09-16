// show/showTables.go
// pgtools
// Rewritten: 2025/09/15
// Lists user tables for one or more DBs (provided as args or types.DefaultDB).
// No exclude list, no size-sorting flag. Always sorted by schema, DB, table.

package show

import (
	"context"
	"fmt"
	"os"
	"pgtools/types"
	"sort"
	"strings"

	"pgtools/shared"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowTables connects to each DB in `dbs` using env-based DSNs and prints a single
// pretty table. Rows are always sorted by Schema, DB, Table (case-insensitive).
func ShowTables(ctx context.Context, cfg *types.DBConfig, dbs []string) *ce.CustomError {
	var rowsAll []TableRow

	for _, dbname := range dbs {
		// Build a DSN for this DB and open a short-lived pool
		dsn := shared.BuildDSN(cfg, dbname)
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return &ce.CustomError{Code: 801, Title: "Error creating DB connection", Message: fmt.Sprintf("%s: %v", dbname, err)}
		}

		rs, ferr := fetchTables(ctx, pool)
		pool.Close()
		if ferr != nil {
			return ferr
		}
		for i := range rs {
			rs[i].DB = dbname
		}
		rowsAll = append(rowsAll, rs...)
	}

	// Always sort: schemaName, DBName, TableName
	sort.Slice(rowsAll, func(i, j int) bool {
		si, sj := strings.ToLower(rowsAll[i].Schema), strings.ToLower(rowsAll[j].Schema)
		if si != sj {
			return si < sj
		}
		di, dj := strings.ToLower(rowsAll[i].DB), strings.ToLower(rowsAll[j].DB)
		if di != dj {
			return di < dj
		}
		return strings.ToLower(rowsAll[i].Table) < strings.ToLower(rowsAll[j].Table)
	})

	// Render
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{"DB", "SCHEMA", "TABLE", "OWNER", "ROWS(EST)", "TOTAL", "TABLE", "INDEX", "TOAST", "PK"})
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.DrawBorder = true
	tw.Style().Format.Header = text.FormatDefault

	if PagedOutput {
		_, terminalHeight := hf.GetTerminalSize()
		if terminalHeight == 0 {
			return &ce.CustomError{Title: "Unable to render table", Message: "Unable to get terminal size", Code: 811}
		}
		tw.SetPageSize(terminalHeight / 3 * 2) // pager is set to 2/3 of the terminal height
	}

	for _, r := range rowsAll {
		tw.AppendRow(table.Row{
			r.DB,
			r.Schema,
			r.Table,
			r.Owner,
			r.RowsEst,
			shared.HumanizeBytes(r.TotalB),
			shared.HumanizeBytes(r.TableB),
			shared.HumanizeBytes(r.IndexB),
			shared.HumanizeBytes(r.ToastB),
			boolToYesNo(r.HasPK),
		})
	}
	tw.SortBy([]table.SortBy{
		{Name: "DB", Mode: table.Asc},
		{Name: "ROWS(EST)", Mode: table.DscNumeric},
	})
	tw.Render()
	return nil
}

func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// fetchTables returns all user tables with size metrics for the *current* database.
func fetchTables(ctx context.Context, pool *pgxpool.Pool) ([]TableRow, *ce.CustomError) {
	const q = `
WITH rel AS (
  SELECT
    n.nspname AS schema,
    c.relname AS tbl,
    pg_catalog.pg_get_userbyid(c.relowner) AS owner,
    COALESCE(c.reltuples::bigint, 0) AS rows_est,
    pg_total_relation_size(c.oid) AS total_b,
    pg_relation_size(c.oid) AS table_b,
    COALESCE(pg_indexes_size(c.oid), 0) AS index_b,
    COALESCE(pg_total_relation_size(c.reltoastrelid), 0) AS toast_b,
    EXISTS (
      SELECT 1 FROM pg_index i
      WHERE i.indrelid = c.oid AND i.indisprimary
    ) AS has_pk
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE c.relkind IN ('r','p') -- ordinary tables and partitions
    AND n.nspname NOT IN ('pg_catalog','information_schema','pg_toast')
)
SELECT schema, tbl, owner, rows_est, total_b, table_b, index_b, toast_b, has_pk
FROM rel
ORDER BY schema, tbl;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, &ce.CustomError{Code: 807, Title: "Error querying tables", Message: err.Error()}
	}
	defer rows.Close()

	var items []TableRow
	for rows.Next() {
		var r TableRow
		if err := rows.Scan(&r.Schema, &r.Table, &r.Owner, &r.RowsEst, &r.TotalB, &r.TableB, &r.IndexB, &r.ToastB, &r.HasPK); err != nil {
			return nil, &ce.CustomError{Code: 808, Title: "Error scanning table row", Message: err.Error()}
		}
		items = append(items, r)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 809, Title: "Row iteration error (tables)", Message: rows.Err().Error()}
	}
	return items, nil
}

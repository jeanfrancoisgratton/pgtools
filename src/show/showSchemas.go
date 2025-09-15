// show/showSchemas.go
// pgtools
// Updated: resolve 'pg_database_owner' to the real DB owner
// 2025/09/14

package show

import (
	"context"
	"os"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// SchemaRow is a single display row for the aggregated schemas table.

// CollectSchemas queries one database connection (pool) and returns rows tagged with dbName.
func CollectSchemas(ctx context.Context, pool *pgxpool.Pool, dbName string) ([]SchemaRow, *ce.CustomError) {
	const q = `
WITH dbinfo AS (
  SELECT d.datdba::regrole::text AS db_owner
  FROM pg_database d
  WHERE d.datname = current_database()
),
user_tables AS (
  SELECT schemaname, COUNT(*) AS cnt
  FROM pg_tables
  WHERE schemaname NOT IN ('pg_catalog','information_schema')
  GROUP BY schemaname
),
user_views AS (
  SELECT schemaname, COUNT(*) AS cnt
  FROM pg_views
  WHERE schemaname NOT IN ('pg_catalog','information_schema')
  GROUP BY schemaname
),
schema_sizes AS (
  SELECT n.nspname AS schemaname,
         COALESCE(SUM(pg_total_relation_size(c.oid)), 0) AS size_bytes
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname NOT IN ('pg_catalog','information_schema')
    AND c.relkind IN ('r','p','m','i','S','t')  -- tables, partitions, matviews, indexes, sequences, toast
  GROUP BY n.nspname
)
SELECT n.nspname AS schema,
       CASE WHEN r.rolname = 'pg_database_owner' THEN dbinfo.db_owner ELSE r.rolname END AS owner,
       COALESCE(t.cnt, 0) AS tables,
       COALESCE(v.cnt, 0) AS views,
       pg_size_pretty(COALESCE(s.size_bytes, 0)) AS total_size
FROM pg_namespace n
LEFT JOIN pg_roles r    ON r.oid = n.nspowner
LEFT JOIN user_tables t ON t.schemaname = n.nspname
LEFT JOIN user_views v  ON v.schemaname = n.nspname
LEFT JOIN schema_sizes s ON s.schemaname = n.nspname
CROSS JOIN dbinfo
WHERE n.nspname NOT IN ('pg_catalog','information_schema')
ORDER BY n.nspname;
`

	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, &ce.CustomError{Code: 801, Title: "Error listing schemas", Message: err.Error()}
	}
	defer rows.Close()

	out := make([]SchemaRow, 0, 32)
	for rows.Next() {
		var r SchemaRow
		r.DB = dbName
		if err := rows.Scan(&r.Schema, &r.Owner, &r.Tables, &r.Views, &r.TotalSize); err != nil {
			return nil, &ce.CustomError{Code: 802, Title: "Error scanning schemas", Message: err.Error()}
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, &ce.CustomError{Code: 803, Title: "Error scanning rows", Message: err.Error()}
	}

	return out, nil
}

// RenderSchemas prints one aggregated table for all DBs.
func RenderSchemas(rows []SchemaRow) {
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.SetStyle(table.StyleRounded)
	tw.Style().Format.Header = text.FormatUpper

	tw.AppendHeader(table.Row{"DB", "Schema", "Owner", "Tables", "Views", "Total Size"})

	// Sort by DB, then Schema for stable, predictable output
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].DB == rows[j].DB {
			return rows[i].Schema < rows[j].Schema
		}
		return rows[i].DB < rows[j].DB
	})

	for _, r := range rows {
		tw.AppendRow(table.Row{r.DB, r.Schema, r.Owner, r.Tables, r.Views, r.TotalSize})
	}
	tw.Render()
}

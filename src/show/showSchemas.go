// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 11:15
// Original filename: src/show/showSchemas.go

package show

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListSchemas prints all schemas with owner, #tables, #views, and total size.
func ShowSchemas(ctx context.Context, pool *pgxpool.Pool) error {
	const q = `
WITH rels AS (
  SELECT n.oid AS nspoid,
         n.nspname,
         n.nspowner,
         c.oid AS reloid,
         c.relkind
  FROM pg_namespace n
  LEFT JOIN pg_class c ON c.relnamespace = n.oid
  WHERE n.nspname NOT LIKE 'pg_%'
    AND n.nspname <> 'information_schema'
),
agg AS (
  SELECT nspoid,
         COUNT(*) FILTER (WHERE relkind = 'r') AS tables,
         COUNT(*) FILTER (WHERE relkind IN ('v','m')) AS views,
         COALESCE(SUM(pg_total_relation_size(reloid)), 0) AS total_bytes
  FROM rels
  GROUP BY nspoid
)
SELECT r.nspname AS schema,
       pg_get_userbyid(r.nspowner) AS owner,
       COALESCE(a.tables, 0) AS tables,
       COALESCE(a.views, 0) AS views,
       pg_size_pretty(COALESCE(a.total_bytes, 0)) AS total_size
FROM (SELECT DISTINCT nspoid, nspname, nspowner FROM rels) r
LEFT JOIN agg a USING (nspoid)
ORDER BY r.nspname;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return fmt.Errorf("query schemas: %w", err)
	}
	defer rows.Close()

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{"SCHEMA", "OWNER", "TABLES", "VIEWS", "TOTAL SIZE"})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "TABLES", Align: text.AlignRight},
		{Name: "VIEWS", Align: text.AlignRight},
	})

	for rows.Next() {
		var schema, owner, totalSize string
		var tablesCnt, viewsCnt int64
		if err := rows.Scan(&schema, &owner, &tablesCnt, &viewsCnt, &totalSize); err != nil {
			return fmt.Errorf("scan schema: %w", err)
		}
		tw.AppendRow(table.Row{schema, owner, tablesCnt, viewsCnt, totalSize})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows (schemas): %w", err)
	}

	tw.Render()
	return nil
}

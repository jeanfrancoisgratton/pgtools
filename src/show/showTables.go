// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 11:16
// Original filename: src/show/showTables.go

package show

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListTables prints tables with schema, owner, est rows, sizes, and PK presence.
func ShowTables(ctx context.Context, pool *pgxpool.Pool) error {
	const q = `
WITH user_tables AS (
  SELECT
    n.nspname AS schema,
    c.relname AS table,
    c.oid     AS reloid,
    c.relowner AS relowner
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE c.relkind = 'r'
    AND n.nspname NOT LIKE 'pg_%'
    AND n.nspname <> 'information_schema'
),
sizes AS (
  SELECT
    reloid,
    pg_total_relation_size(reloid) AS total_bytes,
    pg_relation_size(reloid)       AS table_bytes,
    pg_indexes_size(reloid)        AS index_bytes,
    COALESCE(pg_total_relation_size(reloid) - pg_relation_size(reloid) - pg_indexes_size(reloid), 0) AS toast_bytes
  FROM user_tables
),
rows_est AS (
  SELECT s.relid AS reloid,
         s.n_live_tup::bigint AS live_rows
  FROM pg_stat_all_tables s
),
pk AS (
  SELECT c.oid AS reloid,
         EXISTS (
           SELECT 1
           FROM pg_index i
           WHERE i.indrelid = c.oid AND i.indisprimary
         ) AS has_pk
  FROM pg_class c
)
SELECT
  t.schema,
  t.table,
  pg_get_userbyid(t.relowner) AS owner,
  COALESCE(r.live_rows, 0)    AS rows,
  pg_size_pretty(s.total_bytes) AS total_size,
  pg_size_pretty(s.table_bytes) AS table_size,
  pg_size_pretty(s.index_bytes) AS index_size,
  pg_size_pretty(s.toast_bytes) AS toast_size,
  p.has_pk
FROM user_tables t
LEFT JOIN sizes s USING (reloid)
LEFT JOIN rows_est r USING (reloid)
LEFT JOIN pk p USING (reloid)
ORDER BY t.schema, t.table;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return fmt.Errorf("query tables: %w", err)
	}
	defer rows.Close()

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{
		"SCHEMA", "TABLE", "OWNER", "ROWS(EST)", "TOTAL", "TABLE", "INDEX", "TOAST", "PK",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ROWS(EST)", Align: text.AlignRight},
		{Name: "TOTAL", Align: text.AlignRight},
		{Name: "TABLE", Align: text.AlignRight},
		{Name: "INDEX", Align: text.AlignRight},
		{Name: "TOAST", Align: text.AlignRight},
		{Name: "PK", Align: text.AlignCenter},
	})

	for rows.Next() {
		var (
			schema, tbl, owner                 string
			totalSz, tableSz, indexSz, toastSz string
			rowsEst                            int64
			hasPK                              bool
		)
		if err := rows.Scan(&schema, &tbl, &owner, &rowsEst, &totalSz, &tableSz, &indexSz, &toastSz, &hasPK); err != nil {
			return fmt.Errorf("scan table: %w", err)
		}
		tw.AppendRow(table.Row{
			schema, tbl, owner, rowsEst, totalSz, tableSz, indexSz, toastSz, hasPK,
		})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows (tables): %w", err)
	}

	tw.Render()
	return nil
}

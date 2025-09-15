// src/conf/list.go
// pgtools
// Lists all PostgreSQL configuration parameters (equivalent to SHOW ALL)
// and renders a pretty table similar to environment.ListEnvironments().

package conf

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// CollectAll queries pg_settings and returns all rows.
func CollectAll(ctx context.Context, pool *pgxpool.Pool) ([]Row, *ce.CustomError) {
	const q = `
SELECT name, setting, COALESCE(unit,''), source, category, short_desc
FROM pg_settings
ORDER BY name;
`
	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, &ce.CustomError{Code: 901, Title: "Error querying pg_settings", Message: err.Error()}
	}
	defer rows.Close()

	var out []Row
	for rows.Next() {
		var r Row
		if sErr := rows.Scan(&r.Name, &r.Setting, &r.Unit, &r.Source, &r.Category, &r.ShortDesc); sErr != nil {
			return nil, &ce.CustomError{Code: 902, Title: "Error scanning pg_settings", Message: sErr.Error()}
		}
		out = append(out, r)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 903, Title: "Error iterating pg_settings", Message: rows.Err().Error()}
	}
	return out, nil
}

// CollectByNames fetches a subset of settings by their names (case-insensitive).
func CollectByNames(ctx context.Context, pool *pgxpool.Pool, names []string) ([]Row, *ce.CustomError) {
	if len(names) == 0 {
		return []Row{}, nil
	}
	const q = `
SELECT name, setting, COALESCE(unit,''), source, category, short_desc
FROM pg_settings
WHERE lower(name) = ANY($1::text[])
ORDER BY name;
`
	lowered := make([]string, 0, len(names))
	for _, n := range names {
		lowered = append(lowered, strings.ToLower(n))
	}

	rows, err := pool.Query(ctx, q, lowered)
	if err != nil {
		return nil, &ce.CustomError{Code: 904, Title: "Error querying pg_settings subset", Message: err.Error()}
	}
	defer rows.Close()

	var out []Row
	for rows.Next() {
		var r Row
		if sErr := rows.Scan(&r.Name, &r.Setting, &r.Unit, &r.Source, &r.Category, &r.ShortDesc); sErr != nil {
			return nil, &ce.CustomError{Code: 905, Title: "Error scanning pg_settings subset", Message: sErr.Error()}
		}
		out = append(out, r)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 906, Title: "Error iterating pg_settings subset", Message: rows.Err().Error()}
	}
	return out, nil
}

// ellipsize truncates to max characters (rune-safe) and appends "..." if needed.
func ellipsize(s string, max int) string {
	if max <= 0 {
		return ""
	}
	rs := []rune(strings.TrimSpace(s))
	if len(rs) <= max {
		return string(rs)
	}
	if max <= 3 {
		return strings.Repeat(".", max)
	}
	return string(rs[:max-3]) + "..."
}

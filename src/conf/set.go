// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/15 07:52
// Original filename: src/conf/set.go

package conf

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// Set sets a server-wide parameter using ALTER SYSTEM and reloads the config.
// Requires superuser or appropriate privileges. The change is persisted to
// postgresql.auto.conf.
func Set(ctx context.Context, pool *pgxpool.Pool, name, value string) *ce.CustomError {
	if !ValidName.MatchString(name) {
		return &ce.CustomError{Code: 911, Title: "Invalid setting name", Message: "Only letters/digits/underscore/dot are allowed"}
	}

	// Use parameterized literal for value; name cannot be parameterized as an identifier.
	_, err := pool.Exec(ctx, "ALTER SYSTEM SET "+name+" = $1", value)
	if err != nil {
		return &ce.CustomError{Code: 912, Title: "ALTER SYSTEM failed", Message: err.Error()}
	}

	// Reload to apply ASAP where possible.
	if _, rerr := pool.Exec(ctx, "SELECT pg_reload_conf()"); rerr != nil {
		return &ce.CustomError{Code: 913, Title: "pg_reload_conf() failed", Message: rerr.Error()}
	}
	return nil
}

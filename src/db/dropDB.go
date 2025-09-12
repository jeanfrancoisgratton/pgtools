// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 04:05
// Original filename: src/db/drop.go

package db

import (
	"context"
	"pgtools/shared"

	ce "github.com/jeanfrancoisgratton/customError/v2"
	"pgtools/logging"
	"pgtools/types"
)

// DropDatabase drops a database. If force=true, disconnect sessions first.
// Implementation tries:
//  1. DROP DATABASE <name> WITH (FORCE)
//  2. if (1) fails (older PG), terminate backends then plain DROP DATABASE
func DropDatabase(cfg *types.DBConfig, dbname string, force bool) *ce.CustomError {
	logging.Debugf("Entering function: DropDatabase(%s, force=%v)", dbname, force)

	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// Fast path: try native FORCE if requested (PG >= 13)
	if force {
		stmt := "DROP DATABASE " + shared.QuoteIdent(dbname) + " WITH (FORCE)"
		if _, e := conn.Exec(context.Background(), stmt); e == nil {
			logging.Infof("Dropped database %s using WITH (FORCE)", dbname)
			return nil
		}
		// Fall through to manual terminate + drop
	}

	// If force requested (or previous attempt failed), terminate backends and drop.
	if force {
		kill := `
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = $1
  AND pid <> pg_backend_pid();`
		if _, e := conn.Exec(context.Background(), kill, dbname); e != nil {
			return &ce.CustomError{Code: 711, Title: "Terminate backends failed", Message: e.Error()}
		}
	}

	// Plain drop (also the non-force path)
	stmt := "DROP DATABASE " + shared.QuoteIdent(dbname)
	if _, e := conn.Exec(context.Background(), stmt); e != nil {
		return &ce.CustomError{Code: 712, Title: "DROP DATABASE failed", Message: e.Error()}
	}

	logging.Infof("Dropped database %s", dbname)
	return nil
}

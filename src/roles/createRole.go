// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:42
// Original filename: src/roles/createRole.go

package roles

import (
	"context"
	"fmt"
	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// CreateRole creates a new role with provided attributes.
func CreateRole(cfg *types.DBConfig, name string, opts types.RoleOptions) *ce.CustomError {
	logging.Debugf("Entering function: CreateRole(%s)", name)

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	attrs := make([]string, 0, 10)

	if opts.Superuser != nil {
		if *opts.Superuser {
			attrs = append(attrs, "SUPERUSER")
		} else {
			attrs = append(attrs, "NOSUPERUSER")
		}
	}
	if opts.CreateDB != nil {
		if *opts.CreateDB {
			attrs = append(attrs, "CREATEDB")
		} else {
			attrs = append(attrs, "NOCREATEDB")
		}
	}
	if opts.CreateRole != nil {
		if *opts.CreateRole {
			attrs = append(attrs, "CREATEROLE")
		} else {
			attrs = append(attrs, "NOCREATEROLE")
		}
	}
	if opts.Inherit != nil {
		if *opts.Inherit {
			attrs = append(attrs, "INHERIT")
		} else {
			attrs = append(attrs, "NOINHERIT")
		}
	}
	if opts.Login != nil {
		if *opts.Login {
			attrs = append(attrs, "LOGIN")
		} else {
			attrs = append(attrs, "NOLOGIN")
		}
	}
	if opts.Replication != nil {
		if *opts.Replication {
			attrs = append(attrs, "REPLICATION")
		} else {
			attrs = append(attrs, "NOREPLICATION")
		}
	}
	if opts.BypassRLS != nil {
		if *opts.BypassRLS {
			attrs = append(attrs, "BYPASSRLS")
		} else {
			attrs = append(attrs, "NOBYPASSRLS")
		}
	}

	q := fmt.Sprintf("CREATE ROLE %s", quoteIdent(name))
	if len(attrs) > 0 {
		q += " WITH " + strings.Join(attrs, " ")
	}

	// CREATE ROLE ... PASSWORD '<value>' (if provided)
	if opts.Password != nil && !opts.ClearPassword {
		q += " PASSWORD " + "'" + strings.ReplaceAll(*opts.Password, `'`, `''`) + "'"
	}

	if _, execErr := conn.Exec(context.Background(), q); execErr != nil {
		return &ce.CustomError{Code: 310, Title: "CREATE ROLE failed", Message: execErr.Error()}
	}

	logging.Infof("Created role: %s", name)
	return nil
}

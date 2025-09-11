// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:47
// Original filename: src/roles/alterRole.go

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

// AlterRole changes role attributes. Only non-nil options are applied.
func AlterRole(cfg *types.DBConfig, name string, opts types.RoleOptions) *ce.CustomError {
	logging.Debugf("Entering function: AlterRole(%s)", name)

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

	// Apply toggles in one ALTER ... WITH statement if any are present.
	if len(attrs) > 0 {
		q := fmt.Sprintf("ALTER ROLE %s WITH %s", quoteIdent(name), strings.Join(attrs, " "))
		if _, execErr := conn.Exec(context.Background(), q); execErr != nil {
			return &ce.CustomError{Code: 320, Title: "ALTER ROLE (attributes) failed", Message: execErr.Error()}
		}
	}

	// Password handling (independent statements for precise errors)
	if opts.ClearPassword {
		q := fmt.Sprintf("ALTER ROLE %s PASSWORD NULL", quoteIdent(name))
		if _, execErr := conn.Exec(context.Background(), q); execErr != nil {
			return &ce.CustomError{Code: 321, Title: "ALTER ROLE (password NULL) failed", Message: execErr.Error()}
		}
	} else if opts.Password != nil {
		q := fmt.Sprintf("ALTER ROLE %s PASSWORD $1", quoteIdent(name))
		if _, execErr := conn.Exec(context.Background(), q, *opts.Password); execErr != nil {
			return &ce.CustomError{Code: 322, Title: "ALTER ROLE (password) failed", Message: execErr.Error()}
		}
	}

	logging.Infof("Altered role: %s", name)
	return nil
}

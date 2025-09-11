// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:50
// Original filename: src/roles/password.go

package roles

import (
	"context"
	"fmt"
	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// ChangePassword changes or clears a role's password.
// pw == nil && clear==false  => no-op (nothing to do)
// clear == true              => PASSWORD NULL
// pw != nil                  => set to value
func ChangePassword(cfg *types.DBConfig, name string, pw *string, clear bool) *ce.CustomError {
	logging.Debugf("Entering function: ChangePassword(%s)", name)

	if pw == nil && !clear {
		return nil
	}

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	if clear {
		stmt := fmt.Sprintf("ALTER ROLE %s PASSWORD NULL", quoteIdent(name))
		if _, execErr := conn.Exec(context.Background(), stmt); execErr != nil {
			return &ce.CustomError{Code: 340, Title: "ALTER ROLE (password NULL) failed", Message: execErr.Error()}
		}
	} else {
		stmt := fmt.Sprintf("ALTER ROLE %s PASSWORD $1", quoteIdent(name))
		if _, execErr := conn.Exec(context.Background(), stmt, *pw); execErr != nil {
			return &ce.CustomError{Code: 341, Title: "ALTER ROLE (password) failed", Message: execErr.Error()}
		}
	}

	logging.Infof("Changed password for role: %s", name)
	return nil
}

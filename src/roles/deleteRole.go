// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:47
// Original filename: src/roles/deleteRole.go

package roles

import (
	"context"
	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// DropRole drops a role. If cascade is true, append CASCADE. If ifExists is true, prepend IF EXISTS.
func DropRole(cfg *types.DBConfig, name string, cascade, ifExists bool) *ce.CustomError {
	logging.Debugf("Entering function: DropRole(%s)", name)

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	stmt := "DROP ROLE "
	if ifExists {
		stmt += "IF EXISTS "
	}
	stmt += quoteIdent(name)
	if cascade {
		stmt += " CASCADE"
	}

	if _, execErr := conn.Exec(context.Background(), stmt); execErr != nil {
		return &ce.CustomError{Code: 330, Title: "DROP ROLE failed", Message: execErr.Error()}
	}

	logging.Infof("Dropped role: %s", name)
	return nil
}

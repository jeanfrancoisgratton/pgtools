// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 04:06
// Original filename: src/db/createDB.go

package db

import (
	"context"
	"pgtools/shared"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"
	"pgtools/logging"
	"pgtools/types"
)

// CreateDatabase creates a new empty database. If owner is non-empty, sets OWNER.
func CreateDatabase(cfg *types.DBConfig, dbname string, owner string) *ce.CustomError {
	logging.Debugf("Entering function: CreateDatabase(%s, owner=%q)", dbname, owner)

	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	stmt := "CREATE DATABASE " + shared.QuoteIdent(dbname)
	if strings.TrimSpace(owner) != "" {
		stmt += " OWNER " + shared.QuoteIdent(owner)
	}

	if _, e := conn.Exec(context.Background(), stmt); e != nil {
		return &ce.CustomError{Code: 700, Title: "CREATE DATABASE failed", Message: e.Error()}
	}

	logging.Infof("Created database %s (owner=%q)", dbname, owner)
	return nil
}

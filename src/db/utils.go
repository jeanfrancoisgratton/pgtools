// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/09 05:21
// Original filename: src/db/utils.go

package db

import (
	"context"
	"database/sql"
	"fmt"
	"pgtools/logging"
	"pgtools/shared"
	"pgtools/types"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// this properly closes the connection where the connection might have been closed before conn.Close() was called
func safeClose(conn *pgx.Conn) {
	if conn == nil {
		return
	}
	defer func() {
		recover() // suppress panics on buggy conn.pgConn
	}()
	_ = conn.Close(context.Background())
}

func Connect(cfg *types.DBConfig, database string) (*pgx.Conn, *ce.CustomError) {
	logging.Debugf("Entering function: Connect")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := shared.BuildDSN(cfg, database)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, &ce.CustomError{Title: "Connection failure", Message: err.Error()}
	}
	return conn, nil
}

func safeSQLValue(field sql.NullString, fallback string, singleQuotes bool) string {
	quote := ""
	if singleQuotes {
		quote = "'"
	}
	if field.Valid {
		return fmt.Sprintf("%s%s%s", quote, field.String, quote)
	}
	return fmt.Sprintf("%s%s%s", quote, fallback, quote)
}

func stripChars(s string, prfx string) string {
	if strings.HasPrefix(s, prfx) && strings.HasSuffix(s, prfx) {
		return s[1 : len(s)-1]
	} else {
		return s
	}
}

// ShowVersion connects to the PostgreSQL server and returns the output of `SELECT version()`
func ShowDBServerVersion(cfg *types.DBConfig) (string, *ce.CustomError) {
	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var version string
	nerr := conn.QueryRow(context.Background(), `SHOW server_version`).Scan(&version)
	if nerr != nil {
		return "", &ce.CustomError{Code: 204, Title: "Failed to query server_version", Message: nerr.Error()}
	}

	return version, nil
}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Rewritten: 2025/09/12
// Original filename: src/db/backup.go
//
// Purpose:
//   - Remove dependency on pgtools/show to avoid package cycles.
//   - Keep behavior: determine DB show (respecting -a), create archive (.sql[.gz]),
//     and dump each database by calling writeDatabaseSQL().
//

package db

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"strings"

	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// BackupDatabase dumps one or more databases into a single SQL file.
// Usage semantics (from cmd layer):
//
//	pgtools db backup [-a] db1 [db2 ...] archive_name
//
// The last argument is always the target archive filename. If it ends with .gz,
// output is gzip-compressed. The base SQL always uses a .sql extension.
func BackupDatabase(cfg *types.DBConfig, inOutArgs []string) *ce.CustomError {
	logging.Debugf("Entering function: db.BackupDatabase")

	if len(inOutArgs) < 1 {
		return &ce.CustomError{Code: 90, Title: "Invalid arguments", Message: "missing archive filename"}
	}

	// Archive filename is the last argument
	archive := inOutArgs[len(inOutArgs)-1]

	// Handle filename extensions: ensure .sql, optional .gz
	gzExt := false
	if strings.HasSuffix(archive, ".gz") {
		gzExt = true
		archive = strings.TrimSuffix(archive, ".gz")
	}
	if strings.HasSuffix(archive, ".sql") {
		archive = strings.TrimSuffix(archive, ".sql")
	}
	archive += ".sql"
	if gzExt {
		archive += ".gz"
	}

	// Build database show
	var dbnames []string
	if types.AllDBs {
		var cerr *ce.CustomError
		if dbnames, cerr = getDatabaseNames(cfg); cerr != nil {
			logging.Errorf("Error code %d -> %s : %s", cerr.Code, cerr.Title, cerr.Message)
			return cerr
		}
	} else {
		if len(inOutArgs) < 2 {
			return &ce.CustomError{Code: 91, Title: "Invalid arguments", Message: "no databases specified"}
		}
		dbnames = inOutArgs[:len(inOutArgs)-1]
	}

	// Open output file
	file, err := os.Create(archive)
	if err != nil {
		return &ce.CustomError{Code: 92, Title: "Cannot create archive", Message: err.Error()}
	}
	defer func() { _ = file.Close() }()

	var writer io.Writer = file
	var gzWriter *gzip.Writer
	if gzExt {
		gzWriter = gzip.NewWriter(file)
		writer = gzWriter
		defer func() { _ = gzWriter.Close() }()
	}

	// Dump each database
	for _, dbname := range dbnames {
		if err := writeDatabaseSQL(cfg, dbname, writer); err != nil {
			logging.Errorf("Error code %d -> %s : %s", err.Code, err.Title, err.Message)
			return err
		}
	}

	return nil
}

// getDatabaseNames returns all non-template database names, ordered by name.
func getDatabaseNames(cfg *types.DBConfig) ([]string, *ce.CustomError) {
	logging.Debugf("Entering function: db.getDatabaseNames")

	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, qerr := conn.Query(context.Background(), `
        SELECT datname
        FROM pg_database
        WHERE datistemplate = false
        ORDER BY datname`)
	if qerr != nil {
		return nil, &ce.CustomError{Code: 101, Title: "Unable to show databases", Message: qerr.Error()}
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, &ce.CustomError{Code: 102, Title: "Scan error", Message: err.Error()}
		}
		names = append(names, n)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 103, Title: "List iteration failed", Message: rows.Err().Error()}
	}

	return names, nil
}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/08 17:05
// Original filename: src/db/backup.go

package db

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"

	"pgtools/logging"
	"pgtools/types"
)

func BackupDatabase(cfg *types.DBConfig, inOutArgs []string) *ce.CustomError {
	logging.Debugf("Entering function: BackupDatabase")

	if len(inOutArgs) == 0 {
		logging.Errorf("Error code 201 -> Missing archive name")
		os.Exit(201)
	}

	archive := inOutArgs[len(inOutArgs)-1]
	// archive filename extension handling
	// full filename should always be $filename.sql[.gz]
	var gzExt bool
	if gzExt = strings.HasSuffix(archive, ".gz"); gzExt {
		archive = strings.TrimSuffix(archive, ".gz")
	}
	if strings.HasSuffix(archive, ".sql") {
		archive = strings.TrimSuffix(archive, ".sql")
	}
	archive += ".sql"
	if gzExt {
		archive += ".gz"
	}

	file, err := os.Create(archive)
	if err != nil {
		logging.Errorf("Error code 202 -> Unable to create archive: %v", err)
		return &ce.CustomError{Code: 202, Title: "Failed to create archive", Message: err.Error()}
	}
	defer file.Close()

	var outputWriter io.Writer
	outputWriter = file

	var gzipWriter *gzip.Writer
	if strings.HasSuffix(archive, ".gz") {
		gzipWriter = gzip.NewWriter(file)
		outputWriter = gzipWriter
		defer gzipWriter.Close()
	}

	// Handle -u flag
	if types.UserRoles {
		if err := DumpGlobalRoles(cfg, outputWriter); err != nil {
			logging.Errorf("Error code %d -> %s : %s", err.Code, err.Title, err.Message)
			return err
		}
		return nil
	}

	// Determine database list
	var dbnames []string
	if types.AllDBs {
		var cerr *ce.CustomError
		if dbnames, cerr = ListDatabases(cfg); cerr != nil {
			logging.Errorf("Error code %d -> %s : %s", cerr.Code, cerr.Title, cerr.Message)
			return cerr
		}
	} else {
		if len(inOutArgs) < 2 {
			logging.Errorf(fmt.Sprintf("Error code 203 -> Missing database name(s)"))
			os.Exit(203)
		}
		dbnames = inOutArgs[:len(inOutArgs)-1]
	}

	// Dump each database
	for _, db := range dbnames {
		if err := writeDatabaseSQL(cfg, db, outputWriter); err != nil {
			logging.Errorf("Error code %d -> %s : %s", err.Code, err.Title, err.Message)
			return err
		}
	}

	return nil
}

// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/08 17:09
// Original filename: src/db/restore.go

package db

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"io"
	"os"
	"pgtool/types"
	"strings"
)

func RestoreDatabase(cfg *types.DBConfig, inOutArgs []string) *ce.CustomError {
	for _, arg := range inOutArgs {
		if err := restoreDB(cfg, arg); err != nil {
			return err
		}
	}
	return nil
}

func restoreDB(cfg *types.DBConfig, arcname string) *ce.CustomError {
	file, err := os.Open(arcname)
	if err != nil {
		return &ce.CustomError{Title: "could not open file", Message: err.Error(), Code: 201}
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(arcname, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return &ce.CustomError{Title: "gzip decompression failed", Message: err.Error(), Code: 202}
		}
		defer gzReader.Close()
		reader = gzReader
	}

	conn, cerr := Connect(cfg, "postgres")
	if cerr != nil {
		return cerr
	}
	defer conn.Close(context.Background())

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 20*1024*1024) // allow up to 20MB lines

	var sb strings.Builder
	var pendingDBName string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		// Handle \c or \connect — store it, delay actual reconnection
		if strings.HasPrefix(trimmed, `\c `) || strings.HasPrefix(trimmed, `\connect `) {
			fields := strings.Fields(trimmed)
			if len(fields) >= 2 {
				pendingDBName = stripChars(fields[1], "\"")
			}
			continue
		}

		sb.WriteString(line)
		sb.WriteString(" ")

		if strings.HasSuffix(trimmed, ";") {
			query := sb.String()
			sb.Reset()

			if types.UserRoles {
				if !strings.Contains(query, "pg_roles") &&
					!strings.Contains(query, "pg_auth_members") &&
					!strings.Contains(query, "pg_shadow") &&
					!strings.Contains(query, "pg_user") {
					continue
				}
			}

			_, err := conn.Exec(context.Background(), query)
			if err != nil {
				return &ce.CustomError{
					Title:   fmt.Sprintf("query execution failed\n%s", query),
					Message: err.Error(),
					Code:    204,
				}
			}

			// After successful execution of the current statement,
			// apply any pending \connect line.
			if pendingDBName != "" {
				conn.Close(context.Background())
				conn, cerr = Connect(cfg, pendingDBName)
				if cerr != nil {
					return cerr
				}
				pendingDBName = ""
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return &ce.CustomError{
			Title:   "error reading standard input",
			Message: err.Error(),
			Code:    205,
		}
	}
	return nil
}

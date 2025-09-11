// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/14 18:00
// Original filename: src/db/dump_data.go

package db

import (
	"context"
	"fmt"
	"io"
	"pgtools/logging"
	"pgtools/types"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

func writeDatabaseSQL(cfg *types.DBConfig, dbName string, writer io.Writer) *ce.CustomError {
	logging.Debugf("Entering writeDatabaseSQL for %s", dbName)

	logging.Infof("Processing database: %s", dbName)
	fmt.Fprintf(writer, "-- Dump of database: %s\n", dbName)
	fmt.Fprintf(writer, "-- Dump created with pgtools on %s\n\n", time.Now().Format("2006.01.02 15:04:05"))
	fmt.Fprintf(writer, "DROP DATABASE IF EXISTS %s;\n", stripChars(dbName, "\""))
	if dbDef, nerr := GetDatabaseDefinition(cfg, dbName); nerr == nil {
		fmt.Fprintf(writer, "%s;\n", dbDef)
	} else {
		logging.Errorf("Code: %v -> %s : %s", nerr.Code, nerr.Title, nerr.Message)
		return nerr
	}
	fmt.Fprintf(writer, "\\c %s\n\n", stripChars(dbName, "\""))

	// Get the DB attributes
	attrs, aerr := GetDBAttributes(dbName, cfg)
	if aerr != nil {
		return aerr
	}
	for _, stmt := range attrs {
		fmt.Fprintf(writer, "%s\n", stmt)
	}
	fmt.Fprintln(writer)

	// Set the DB ownership
	owner, oerr := GetOwnership(dbName, cfg)
	if oerr != nil {
		return oerr
	}
	fmt.Fprintf(writer, "ALTER DATABASE %s OWNER TO %s;\n", stripChars(dbName, "\""), owner)

	// Get the SEQUENCE values
	seqs, serr := GetSequences(cfg, dbName, owner)
	if serr != nil {
		return serr
	}
	for _, s := range seqs {
		_, _ = fmt.Fprintln(writer, s)
	}

	conn, err := Connect(cfg, dbName)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	tables, cerr := ListTables(conn)
	if cerr != nil {
		return cerr
	}

	for _, table := range tables {
		if err := dumpTable(conn, table, writer); err != nil {
			return err
		}
	}

	// Write constraints after all schema and data
	constraints, cerr := GetAllConstraints(dbName, cfg)
	if cerr != nil {
		return cerr
	}
	for _, stmt := range constraints {
		fmt.Fprintf(writer, "%s\n", stmt)
	}
	return nil
}

func dumpTable(conn *pgx.Conn, tableName string, writer io.Writer) *ce.CustomError {
	logging.Debugf("Entering dumpTable for %s", tableName)
	logging.Infof("Dumping table: %s", tableName)

	createSQL, err := buildCreateTableSQL(conn, tableName)
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "\nDROP TABLE IF EXISTS %s;\n%s\n", stripChars(tableName, "\""), createSQL)

	rows, nerr := conn.Query(context.Background(), fmt.Sprintf("SELECT * FROM %s", stripChars(tableName, "\"")))
	if nerr != nil {
		return &ce.CustomError{Code: 105, Title: "Failed to query table", Message: nerr.Error()}
	}
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return &ce.CustomError{Code: 106, Title: "Scan failed", Message: err.Error()}
		}

		fmt.Fprintf(writer, "INSERT INTO %s VALUES (", stripChars(tableName, "\""))
		for i, val := range values {
			if i > 0 {
				fmt.Fprint(writer, ", ")
			}
			if val == nil {
				fmt.Fprint(writer, "NULL")
			} else {
				escaped := strings.ReplaceAll(fmt.Sprint(val), "'", "''")
				fmt.Fprintf(writer, "'%s'", escaped)
			}
		}
		fmt.Fprintln(writer, ");")
	}

	return nil
}

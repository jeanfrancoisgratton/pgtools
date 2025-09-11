// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>

package db

import (
	"context"
	"database/sql"
	"fmt"
	"pgtools/logging"
	"pgtools/types"
	"strings"

	"github.com/jackc/pgx/v5"
	ce "github.com/jeanfrancoisgratton/customError/v2"
)

func buildCreateTableSQL(conn *pgx.Conn, table string) (string, *ce.CustomError) {
	logging.Debugf("Entering function: buildCreateTableSQL")
	ctx := context.Background()

	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`
	rows, err := conn.Query(ctx, query, table)
	if err != nil {
		return "", &ce.CustomError{Code: 108, Title: "Query failed", Message: err.Error()}
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var name, typ, nullable string
		if err := rows.Scan(&name, &typ, &nullable); err != nil {
			return "", &ce.CustomError{Code: 109, Title: "Scan failed", Message: err.Error()}
		}

		line := fmt.Sprintf("%s %s", pgx.Identifier{name}.Sanitize(), typ)
		if nullable == "NO" {
			line += " NOT NULL"
		}
		columns = append(columns, line)
	}
	if err := rows.Err(); err != nil {
		return "", &ce.CustomError{Code: 110, Title: "Rows error", Message: err.Error()}
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (\n    %s\n);", pgx.Identifier{table}.Sanitize(), strings.Join(columns, ",\n    "))
	return createSQL, nil
}

func GetDatabaseDefinition(cfg *types.DBConfig, dbName string) (string, *ce.CustomError) {
	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	query := `SELECT
		d.datname,
		pg_encoding_to_char(d.encoding),
		t.spcname AS tablespace,
		pg_catalog.pg_get_userbyid(d.datdba) AS owner,
		tpl.datname AS template_name,
		d.datcollate,
		d.datctype
	FROM pg_database d
	LEFT JOIN pg_tablespace t ON d.dattablespace = t.oid
	LEFT JOIN pg_database tpl ON d.datistemplate AND d.oid = tpl.oid
	WHERE d.datname = $1;`

	var (
		datname, encoding, owner string
		tablespace               sql.NullString
		templateName             sql.NullString
		collate                  sql.NullString
		ctype                    sql.NullString
	)

	row := conn.QueryRow(context.Background(), query, dbName)
	nerr := row.Scan(&datname, &encoding, &tablespace, &owner, &templateName, &collate, &ctype)
	if nerr != nil {
		return "", &ce.CustomError{Code: 204, Title: "Failed to scan database info", Message: nerr.Error()}
	}

	parts := []string{fmt.Sprintf("CREATE DATABASE %s", stripChars(datname, "\""))}
	parts = append(parts, fmt.Sprintf("WITH TEMPLATE = %s", safeSQLValue(templateName, "template0", false)))
	parts = append(parts, fmt.Sprintf("ENCODING = '%s'", encoding))

	if collate.Valid {
		parts = append(parts, fmt.Sprintf("LC_COLLATE = '%s'", collate.String))
	}
	if ctype.Valid {
		parts = append(parts, fmt.Sprintf("LC_CTYPE = '%s'", ctype.String))
	}
	if tablespace.Valid && tablespace.String != "pg_default" {
		parts = append(parts, fmt.Sprintf("TABLESPACE = \"%s\"", tablespace.String))
	}

	return strings.Join(parts, " ") + ";", nil
}

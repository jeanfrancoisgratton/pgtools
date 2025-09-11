// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/12 15:23
// Original filename: src/db/dbmetadata.go

package db

import (
	"context"
	"database/sql"
	"fmt"
	"pgtools/types"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

func GetOwnership(databaseName string, cfg *types.DBConfig) (string, *ce.CustomError) {
	conn, err := Connect(cfg, databaseName)
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var owner string
	row := conn.QueryRow(context.Background(), `
		SELECT pg_catalog.pg_get_userbyid(datdba)
		FROM pg_catalog.pg_database
		WHERE datname = $1
	`, databaseName)
	if err := row.Scan(&owner); err != nil {
		return "", &ce.CustomError{Code: 401, Title: "Ownership lookup failed", Message: err.Error()}
	}

	return owner, nil
}

func GetDBAttributes(databaseName string, cfg *types.DBConfig) ([]string, *ce.CustomError) {
	conn, err := Connect(cfg, databaseName)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	attributes := []string{
		"statement_timeout",
		"lock_timeout",
		"idle_in_transaction_session_timeout",
		//		"transaction_timeout",
		"client_encoding",
		"standard_conforming_strings",
		"xmloption",
		"client_min_messages",
		"row_security",
	}

	statements := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		var value string
		row := conn.QueryRow(context.Background(), "SHOW "+attr)
		if err := row.Scan(&value); err != nil {
			return nil, &ce.CustomError{Code: 402, Title: "SHOW failed", Message: fmt.Sprintf("SHOW %s: %v", attr, err)}
		}
		// Quote value if it's not purely numeric
		if strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			value = fmt.Sprintf("'%s'", value)
		}
		statements = append(statements, fmt.Sprintf("SET %s = %s;", attr, value))
	}

	return statements, nil
}

func GetSequences(cfg *types.DBConfig, databaseName, dbOwner string) ([]string, *ce.CustomError) {
	conn, err := Connect(cfg, databaseName)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT schemaname, sequencename, start_value, increment_by, min_value, max_value, cache_size, cycle, last_value
FROM pg_sequences
WHERE schemaname = 'public';`

	rows, nerr := conn.Query(context.Background(), query)
	if nerr != nil {
		return nil, &ce.CustomError{Code: 403, Title: "Failed to query sequences", Message: nerr.Error()}
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var schema, name string
		var start, increment, min, max, cache int64
		var isCycled bool
		var last sql.NullInt64

		err := rows.Scan(&schema, &name, &start, &increment, &min, &max, &cache, &isCycled, &last)
		if err != nil {
			return nil, &ce.CustomError{Code: 404, Title: "Failed to scan sequence row", Message: err.Error()}
		}

		fullName := fmt.Sprintf("%s.%s", schema, name)
		stmt := fmt.Sprintf(`CREATE SEQUENCE %s
    INCREMENT BY %d
    %s
    MAXVALUE %d
    %s
    CACHE %d;`,
			fullName,
			increment,
			func() string {
				if min > 0 {
					return fmt.Sprintf("MINVALUE %d", min)
				}
				return "NO MINVALUE"
			}(),
			max,
			func() string {
				if last.Valid {
					return fmt.Sprintf("START WITH %d", last.Int64)
				}
				return ""
			}(),
			cache,
		)

		result = append(result, stmt)
		result = append(result, fmt.Sprintf(`ALTER SEQUENCE %s OWNER TO %s;`, fullName, dbOwner))
	}
	return result, nil
}

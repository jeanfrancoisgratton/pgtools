// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/14 18:18
// Original filename: src/db/constraints.go

package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// GetAllConstraints fetches all constraints and indexes relevant to DB integrity.
func GetAllConstraints(databaseName string, cfg *types.DBConfig) ([]string, *ce.CustomError) {
	conn, err := Connect(cfg, databaseName)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var results []string

	// Primary Keys
	pks, err := fetchPrimaryKeys(conn)
	if err != nil {
		return nil, err
	}
	results = append(results, pks...)

	// Foreign Keys
	fks, err := fetchForeignKeys(conn)
	if err != nil {
		return nil, err
	}
	results = append(results, fks...)

	// Table constraints
	tcs, err := fetchTableConstraints(conn)
	if err != nil {
		return nil, err
	}
	results = append(results, tcs...)

	// Unique Indexes
	uix, err := fetchUniqueIndexes(conn)
	if err != nil {
		return nil, err
	}
	results = append(results, uix...)

	return results, nil
}

func fetchPrimaryKeys(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("fetch primary keys")
	query := `SELECT tc.table_schema, tc.table_name, kcu.column_name, tc.constraint_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		 AND tc.table_schema = kcu.table_schema
		 AND tc.table_name = kcu.table_name
		WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema NOT IN ('pg_catalog', 'information_schema') 
		ORDER BY tc.table_schema, tc.table_name, kcu.ordinal_position;
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, &ce.CustomError{Code: 301, Title: "Primary key query failed", Message: err.Error()}
	}
	defer rows.Close()

	constraints := make(map[string][]string)
	schemaMap := make(map[string]string)

	for rows.Next() {
		var schema, table, column, constraint string
		if err := rows.Scan(&schema, &table, &column, &constraint); err != nil {
			return nil, &ce.CustomError{Code: 302, Title: "Primary key scan failed", Message: err.Error()}
		}
		key := fmt.Sprintf("%s.%s.%s", schema, table, constraint)
		constraints[key] = append(constraints[key], column)
		schemaMap[key] = fmt.Sprintf("%s.%s", schema, table)
	}

	var results []string
	for key, columns := range constraints {
		parts := strings.Split(key, ".")
		constraint := parts[2]
		fullTable := schemaMap[key]
		line := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (%s);",
			fullTable, constraint, strings.Join(columns, ", "))
		results = append(results, line)
	}

	return results, nil
}

func fetchForeignKeys(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("fetch foreign keys")
	query := `
		SELECT
			tc.table_schema,
			tc.table_name,
			tc.constraint_name,
			kcu.column_name,
			ccu.table_schema AS foreign_table_schema,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
		  ON ccu.constraint_name = tc.constraint_name
		WHERE constraint_type = 'FOREIGN KEY' AND tc.table_schema NOT IN ('pg_catalog', 'information_schema');
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, &ce.CustomError{Code: 303, Title: "Foreign key query failed", Message: err.Error()}
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var schema, table, constraint, column, foreignSchema, foreignTable, foreignColumn string
		if err := rows.Scan(&schema, &table, &constraint, &column, &foreignSchema, &foreignTable, &foreignColumn); err != nil {
			return nil, &ce.CustomError{Code: 304, Title: "Foreign key scan failed", Message: err.Error()}
		}
		line := fmt.Sprintf(`ALTER TABLE %s.%s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s(%s);`,
			schema, table, constraint, column, foreignSchema, foreignTable, foreignColumn)
		results = append(results, line)
	}

	return results, nil
}

func fetchTableConstraints(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("fetch table constraints")
	query := `
		SELECT
			tc.table_schema,
			tc.table_name,
			tc.constraint_name,
			tc.constraint_type
		FROM information_schema.table_constraints tc
		WHERE tc.constraint_type NOT IN ('PRIMARY KEY', 'FOREIGN KEY')
		AND tc.table_schema NOT IN ('pg_catalog', 'information_schema');
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, &ce.CustomError{Code: 305, Title: "Table constraint query failed", Message: err.Error()}
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var schema, table, constraint, constraintType string
		if err := rows.Scan(&schema, &table, &constraint, &constraintType); err != nil {
			return nil, &ce.CustomError{Code: 306, Title: "Table constraint scan failed", Message: err.Error()}
		}
		// Placeholder for actual logic; many of these are check constraints
		results = append(results, fmt.Sprintf("-- Skipping unsupported constraint type: %s on %s.%s", constraintType, schema, table))
	}

	return results, nil
}

func fetchUniqueIndexes(conn *pgx.Conn) ([]string, *ce.CustomError) {
	logging.Debugf("fetch unique indexes")
	query := `
		SELECT
			n.nspname AS schema_name,
			t.relname AS table_name,
			i.relname AS index_name,
			string_agg(a.attname, ', ') AS column_names
		FROM pg_class t
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		WHERE ix.indisunique AND NOT ix.indisprimary AND n.nspname NOT IN ('pg_catalog', 'information_schema')
		GROUP BY schema_name, table_name, index_name;
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, &ce.CustomError{Code: 307, Title: "Unique index query failed", Message: err.Error()}
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var schema, table, index, columns string
		if err := rows.Scan(&schema, &table, &index, &columns); err != nil {
			return nil, &ce.CustomError{Code: 308, Title: "Unique index scan failed", Message: err.Error()}
		}
		results = append(results, fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s.%s (%s);", index, schema, table, columns))
	}

	return results, nil
}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/09 18:38
// Original filename: src/db/user_roles.go

package db

import (
	"context"
	"fmt"
	"io"
	"pgtools/logging"
	"strings"

	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

func DumpGlobalRoles(cfg *types.DBConfig, writer io.Writer) *ce.CustomError {
	logging.Debugf("Entering function: DumpGlobalRoles")

	conn, err := Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer safeClose(conn)

	var isSuper bool
	row := conn.QueryRow(context.Background(),
		`SELECT rolsuper FROM pg_roles WHERE rolname = current_user`)
	if err := row.Scan(&isSuper); err != nil {
		return &ce.CustomError{Code: 201, Title: "Privilege check failed", Message: err.Error()}
	}
	if !isSuper {
		return &ce.CustomError{Code: 202, Title: "Insufficient privileges", Message: "Must be superuser to dump roles and users"}
	}

	logging.Infof("Dumping global PostgreSQL roles and users")
	fmt.Fprintln(writer, "-- Global roles and users")

	rows, qerr := conn.Query(context.Background(),
		`SELECT
			rolname,
			rolsuper, rolinherit, rolcreaterole, rolcreatedb,
			rolcanlogin, rolreplication, rolbypassrls,
			rolpassword
		FROM pg_authid
		WHERE rolname !~ '^pg_'
		ORDER BY rolname;
	`)
	if qerr != nil {
		return &ce.CustomError{Code: 203, Title: "Role query failed", Message: qerr.Error()}
	}
	defer rows.Close()

	for rows.Next() {
		var (
			rolname                                                   string
			super, inherit, createrole, createdb, login, repl, bypass bool
			password                                                  *string
		)
		if err := rows.Scan(&rolname, &super, &inherit, &createrole, &createdb,
			&login, &repl, &bypass, &password); err != nil {
			return &ce.CustomError{Code: 204, Title: "Scan failed", Message: err.Error()}
		}

		fmt.Fprintf(writer, "CREATE ROLE %s;\n", rolname)

		attrs := []string{}
		if super {
			attrs = append(attrs, "SUPERUSER")
		} else {
			attrs = append(attrs, "NOSUPERUSER")
		}
		if inherit {
			attrs = append(attrs, "INHERIT")
		} else {
			attrs = append(attrs, "NOINHERIT")
		}
		if createrole {
			attrs = append(attrs, "CREATEROLE")
		} else {
			attrs = append(attrs, "NOCREATEROLE")
		}
		if createdb {
			attrs = append(attrs, "CREATEDB")
		} else {
			attrs = append(attrs, "NOCREATEDB")
		}
		if login {
			attrs = append(attrs, "LOGIN")
		} else {
			attrs = append(attrs, "NOLOGIN")
		}
		if repl {
			attrs = append(attrs, "REPLICATION")
		} else {
			attrs = append(attrs, "NOREPLICATION")
		}
		if bypass {
			attrs = append(attrs, "BYPASSRLS")
		} else {
			attrs = append(attrs, "NOBYPASSRLS")
		}

		if password != nil {
			attrs = append(attrs, fmt.Sprintf("PASSWORD '%s'", *password))
		}

		fmt.Fprintf(writer, "ALTER ROLE %s WITH %s;\n", rolname, strings.Join(attrs, " "))
	}

	if err := rows.Err(); err != nil {
		return &ce.CustomError{Code: 205, Title: "Rows error", Message: err.Error()}
	}

	return nil
}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:06
// Original filename: src/roles/listRoles.go

package roles

import (
	"context"
	"fmt"
	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5" // <-- added: for pgx.Rows
	ce "github.com/jeanfrancoisgratton/customError/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// roleRow holds visible attributes from pg_roles
type roleRow struct {
	Name        string
	Login       bool
	Superuser   bool
	CreateDB    bool
	CreateRole  bool
	Inherit     bool
	Replication bool
	BypassRLS   bool
}

// narrow interface that *pgx.Conn satisfies; avoids inventing db.Conn
type queryConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// ListRoles prints roles. If includeMembers is true, a MEMBERS column is added.
// If verbose is false, we show a compact set of columns; verbose adds all booleans explicitly.
// In quiet mode (types.Quiet), prints one role name per line (plus members if requested).
func ListRoles(cfg *types.DBConfig, includeMembers bool, verbose bool) *ce.CustomError {
	logging.Debugf("Entering function: ListRoles(includeMembers=%v, verbose=%v)", includeMembers, verbose)

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// Fetch roles from pg_roles (visible to all users)
	const qRoles = `
SELECT
  rolname,
  rolcanlogin,
  rolsuper,
  rolcreatedb,
  rolcreaterole,
  rolinherit,
  rolreplication,
  rolbypassrls
FROM pg_catalog.pg_roles
ORDER BY rolname;
`
	rows, qerr := conn.Query(context.Background(), qRoles)
	if qerr != nil {
		return &ce.CustomError{Code: 350, Title: "Failed to query pg_roles", Message: qerr.Error()}
	}
	defer rows.Close()

	var list []roleRow
	for rows.Next() {
		var r roleRow
		if scanErr := rows.Scan(
			&r.Name, &r.Login, &r.Superuser, &r.CreateDB, &r.CreateRole, &r.Inherit, &r.Replication, &r.BypassRLS,
		); scanErr != nil {
			return &ce.CustomError{Code: 351, Title: "Failed to scan pg_roles row", Message: scanErr.Error()}
		}
		list = append(list, r)
	}
	if rows.Err() != nil {
		return &ce.CustomError{Code: 352, Title: "pg_roles iteration error", Message: rows.Err().Error()}
	}

	// Optionally fetch membership map: role -> members[]
	var members map[string][]string
	if includeMembers {
		m, merr := fetchMemberships(conn) // conn satisfies queryConn
		if merr != nil {
			return merr
		}
		members = m
	}

	// Quiet mode: print names (optionally with members)
	if types.Quiet {
		for _, r := range list {
			if includeMembers {
				ms := members[r.Name]
				sort.Strings(ms)
				if len(ms) > 0 {
					fmt.Printf("%s: %s\n", r.Name, strings.Join(ms, ","))
				} else {
					fmt.Println(r.Name)
				}
			} else {
				fmt.Println(r.Name)
			}
		}
		return nil
	}

	// Pretty table
	tw := table.NewWriter()
	tw.SetStyle(table.StyleRounded)
	tw.Style().Options.SeparateRows = false
	tw.Style().Color.Header = text.Colors{text.Bold}
	tw.AppendHeader(buildHeader(includeMembers, verbose))

	for _, r := range list {
		row := buildRow(r, includeMembers, verbose, members)
		tw.AppendRow(row)
	}

	fmt.Println(tw.Render())
	return nil
}

// fetchMemberships builds a map: role -> []member
func fetchMemberships(conn queryConn) (map[string][]string, *ce.CustomError) {
	// Note: pg_auth_members is visible; join to pg_roles for names.
	const q = `
SELECT r.rolname AS role_name, m.rolname AS member_name
FROM pg_catalog.pg_auth_members am
JOIN pg_catalog.pg_roles r ON r.oid = am.roleid
JOIN pg_catalog.pg_roles m ON m.oid = am.member
ORDER BY r.rolname, m.rolname;
`
	rows, err := conn.Query(context.Background(), q)
	if err != nil {
		return nil, &ce.CustomError{Code: 353, Title: "Failed to query role memberships", Message: err.Error()}
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var roleName, memberName string
		if scanErr := rows.Scan(&roleName, &memberName); scanErr != nil {
			return nil, &ce.CustomError{Code: 354, Title: "Failed to scan role memberships row", Message: scanErr.Error()}
		}
		result[roleName] = append(result[roleName], memberName)
	}
	if rows.Err() != nil {
		return nil, &ce.CustomError{Code: 355, Title: "role memberships iteration error", Message: rows.Err().Error()}
	}
	return result, nil
}

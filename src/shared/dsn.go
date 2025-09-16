// src/shared/dsn.go
// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Updated: 2025/09/12

package shared

import (
	"fmt"
	"net/url"
	"pgtools/types"
)

// BuildDSN builds a Postgres connection string from cfg and dbname.
// Example: postgres://user:pass@host:port/dbname?sslmode=verify-full
func BuildDSN(cfg *types.DBConfig, dbname string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   "/" + dbname,
	}

	q := url.Values{}

	// Default sslmode to "prefer" if unset
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "prefer"
	}
	q.Set("sslmode", sslmode)

	// Optional CA bundle to verify server cert; if omitted, system trust store is used
	if cfg.SSLRootCert != "" {
		q.Set("sslrootcert", cfg.SSLRootCert)
	}

	// Reserved for future mutual TLS (client cert auth). Not used yet.
	// Uncommment when ready to use mTLS
	// if cfg.SSLclientCert != "" {
	// 	q.Set("sslcert", cfg.SSLclientCert)
	// }
	// if cfg.SSLclientKey != "" {
	// 	q.Set("sslkey", cfg.SSLclientKey)
	// }

	// Always set application_name unless caller already did
	if q.Get("application_name") == "" {
		q.Set("application_name", types.AppNameKV)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

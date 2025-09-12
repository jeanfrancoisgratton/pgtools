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
// Example: postgres://user:pass@host:port/dbname?sslmode=require
func BuildDSN(cfg *types.DBConfig, dbname string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   "/" + dbname,
	}
	q := url.Values{}
	if cfg.SSLMode != "" {
		q.Set("sslmode", cfg.SSLMode)
	}
	// Optional client cert/key params (ignored if server doesn’t use them)
	if cfg.SSLCert != "" {
		q.Set("sslcert", cfg.SSLCert)
	}
	if cfg.SSLKey != "" {
		q.Set("sslkey", cfg.SSLKey)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

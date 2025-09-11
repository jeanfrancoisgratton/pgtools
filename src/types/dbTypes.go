// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/dbTypes.go

package types

type DBConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	SSLMode     string `json:"sslmode"`
	SSLCert     string `json:"sslcert,omitempty"`
	SSLKey      string `json:"sslkey,omitempty"`
	Description string `json:"comment,omitempty"`
}

//// Rows is a minimal row iterator interface used by higher-level packages (e.g., roles).
//type Rows interface {
//	Next() bool
//	Scan(dest ...any) error
//	Err() error
//	Close()
//}

//// Queryer is the narrow query interface that returns Rows.
//// db/ provides an adapter so callers don't need to import pgx directly.
//type Queryer interface {
//	Query(ctx context.Context, sql string, args ...any) (Rows, error)
//}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/envTypes.go

package types

var CreateOwner string
var DropForce bool

type DBConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	SSLMode       string `json:"sslmode"`
	SSLRootCert   string `json:"sslrootcert,omitempty"`
	SSLclientCert string `json:"sslclientcert,omitempty"`
	SSLclientKey  string `json:"sslclientkey,omitempty"`
	Description   string `json:"comment,omitempty"`
	DefaultDB     string `json:"defaultdb,omitempty"`
}
type DBSize struct {
	Name      string
	SizeBytes int64
}

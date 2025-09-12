// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/dbTypes.go

package types

var CreateOwner string
var DropForce bool
var ListSortBySize bool

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
type DBSize struct {
	Name      string
	SizeBytes int64
}

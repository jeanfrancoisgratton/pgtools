// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/types.go

package types

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
	SSLCert  string `json:"sslcert,omitempty"`
	SSLKey   string `json:"sslkey,omitempty"`
}

var EnvConfigFile = "defaultEnv.json"
var DebugMode = false
var Quiet = false
var AllDBs = false
var UserRoles = false
var LogLevel = "none"

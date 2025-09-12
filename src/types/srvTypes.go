// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 22:50
// Original filename: src/types/srvTypes.go

package types

import "time"

var ReloadLocal bool
var ReloadPIDFile string
var ReloadDataDir string

type LocalControl struct {
	DataDir     string
	PIDFile     string
	StopTimeout time.Duration
}

type ServerInfoStruct struct {
	ServerName string
	ServerPort int
	Version    string
}

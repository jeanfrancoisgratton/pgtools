// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 18:27
// Original filename: src/srv/reload.go

package srv

import (
	"context"
	"syscall"

	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// ReloadRemote uses SQL only (works for local or remote).
func ReloadRemote(cfg *types.DBConfig) *ce.CustomError {
	logging.Debugf("Entering function: ReloadRemote()")

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	if _, e := conn.Exec(context.Background(), "SELECT pg_reload_conf();"); e != nil {
		return &ce.CustomError{Code: 610, Title: "pg_reload_conf() failed", Message: e.Error()}
	}
	logging.Infof("Reloaded config via SQL")
	return nil
}

// ReloadLocal tries SQL first; if that fails, attempts a local SIGHUP using postmaster.pid.
// Use when pgtools is co-located with the server (same host/container).
func ReloadLocal(cfg *types.DBConfig, local types.LocalControl) *ce.CustomError {
	// Prefer SQL — safe and universal
	if err := ReloadRemote(cfg); err == nil {
		return nil
	} else {
		logging.Errorf("SQL reload failed: %s (will attempt SIGHUP if PID is available)", err.Message)
	}

	// Fallback: try to signal the local postmaster
	pid, nerr := resolveLocalPID(local)
	if nerr != nil {
		return nerr
	}
	if e := syscall.Kill(pid, syscall.SIGHUP); e != nil {
		return &ce.CustomError{Code: 611, Title: "SIGHUP failed", Message: e.Error()}
	}
	logging.Infof("Sent SIGHUP to postmaster (pid=%d)", pid)
	return nil
}

// TryBestReload: if 'local' is provided, behave like ReloadLocal; else ReloadRemote.
func TryBestReload(cfg *types.DBConfig, local *types.LocalControl) *ce.CustomError {
	if local != nil {
		return ReloadLocal(cfg, *local)
	}
	return ReloadRemote(cfg)
}

func resolveLocalPID(local types.LocalControl) (int, *ce.CustomError) {
	if local.PIDFile != "" {
		return ReadPostmasterPIDFromFS(local.PIDFile)
	}
	if local.DataDir != "" {
		return ReadPostmasterPIDFromFS(DefaultPIDFile(local.DataDir))
	}
	return 0, &ce.CustomError{Code: 612, Title: "PID resolution failed", Message: "no PIDFile/DataDir available"}
}

// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 22:28
// Original filename: src/srv/srvHelpers.go

package srv

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pgtools/db"
	"pgtools/logging"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// DiscoverDataDirViaSQL returns SHOW data_directory (server-side absolute path).
// Useful for diagnostics or when deciding whether a local fallback is feasible.
func DiscoverDataDirViaSQL(cfg *types.DBConfig) (string, *ce.CustomError) {
	logging.Debugf("Entering function: DiscoverDataDirViaSQL()")

	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var dataDir string
	if e := conn.QueryRow(context.Background(), "SHOW data_directory;").Scan(&dataDir); e != nil {
		return "", &ce.CustomError{Code: 600, Title: "SHOW data_directory failed", Message: e.Error()}
	}
	return dataDir, nil
}

// DefaultPIDFile returns the default postmaster.pid path for a given data directory.
func DefaultPIDFile(dataDir string) string {
	return filepath.Join(dataDir, "postmaster.pid")
}

// ReadPostmasterPIDFromFS reads PID from a local postmaster.pid (first line).
func ReadPostmasterPIDFromFS(pidFile string) (int, *ce.CustomError) {
	logging.Debugf("Entering function: ReadPostmasterPIDFromFS(%s)", pidFile)

	b, e := readAll(pidFile)
	if e != nil {
		return 0, &ce.CustomError{Code: 601, Title: "Read postmaster.pid failed", Message: e.Error()}
	}
	first := firstLine(string(b))
	if first == "" {
		return 0, &ce.CustomError{Code: 602, Title: "Invalid postmaster.pid", Message: "empty first line"}
	}
	pid, conv := strconv.Atoi(first)
	if conv != nil || pid <= 0 {
		return 0, &ce.CustomError{Code: 603, Title: "Invalid postmaster.pid", Message: "bad PID: " + first}
	}
	return pid, nil
}

// helpers

func readAll(path string) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return io.ReadAll(f)
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

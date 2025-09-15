// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 22:36
// Original filename: src/srv/serverVersion.go

package srv

import (
	"context"
	"pgtools/db"
	"pgtools/types"

	ce "github.com/jeanfrancoisgratton/customError/v2"
)

// ShowVersion connects to the PostgreSQL server and returns the output of `SELECT version()`
func ShowDBServerVersion(cfg *types.DBConfig) (*types.ServerInfoStruct, *ce.CustomError) {
	conn, err := db.Connect(cfg, "postgres")
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var version string
	nerr := conn.QueryRow(context.Background(), `SHOW server_version`).Scan(&version)
	if nerr != nil {
		return nil, &ce.CustomError{Code: 204, Title: "Failed to query server_version", Message: nerr.Error()}
	}

	return &types.ServerInfoStruct{ServerName: cfg.Host, ServerPort: cfg.Port, Version: version}, nil
}

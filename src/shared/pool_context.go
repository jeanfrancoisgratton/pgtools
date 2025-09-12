// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 12:32
// Original filename: src/shared/pool_context.go

package shared

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pgtools/environment"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CancellableContext returns a context canceled on SIGINT/SIGTERM, with a sane timeout.
func CancellableContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
	return ctx, cancel
}

// GetPool uses your existing BuildDSN() mechanism to create a pgx pool.
// - Replace cfg.BuildDSN() import/path above to match your codebase.
// - If BuildDSN requires context or args, adjust the call accordingly.
func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := environment.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err.Error())
		os.Exit(err.Code)
	}
	dsn := BuildDSN(cfg, "postgres")

	pc, perr := pgxpool.ParseConfig(dsn)
	if perr != nil {
		return nil, fmt.Errorf("parse DSN: %w", perr)
	}

	pool, perr := pgxpool.NewWithConfig(ctx, pc)
	if perr != nil {
		return nil, fmt.Errorf("open pool: %w", perr)
	}
	return pool, nil
}

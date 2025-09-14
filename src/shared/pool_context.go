// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 12:32
// Original filename: src/shared/pool_context.go

package shared

import (
	"context"
	"os"
	"os/signal"
	"pgtools/environment"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	ce "github.com/jeanfrancoisgratton/customError/v2"
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
func GetPool(ctx context.Context) (*pgxpool.Pool, *ce.CustomError) {
	cfg, err := environment.LoadConfig()
	if err != nil {
		return nil, err
	}
	dsn := BuildDSN(cfg, "postgres")

	pc, perr := pgxpool.ParseConfig(dsn)
	if perr != nil {
		return nil, &ce.CustomError{Code: 803, Title: "Error parsing DSN", Message: perr.Error()}
	}

	pool, perr := pgxpool.NewWithConfig(ctx, pc)
	if perr != nil {
		return nil, &ce.CustomError{Code: 804, Title: "Error connecting to pool", Message: perr.Error()}
	}
	return pool, nil
}

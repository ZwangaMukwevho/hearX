// pkg/storage/mysql.go
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewMySQLConn provides an *sql.DB connected to MySQL,
// retrying Ping() up to 1 minute before erroring out.
func NewMySQLConn(
	lc fx.Lifecycle,
	logger *zap.Logger,
	dsn string,
) (*sql.DB, error) {
	logger.Info("connecting to MySQL", zap.String("dsn", dsn))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("open failed", zap.Error(err))
		return nil, err
	}

	// actively wait up to 60s for the server to be reachable
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if err := db.PingContext(ctx); err == nil {
			logger.Info("MySQL reachable")
			break
		} else {
			logger.Warn("MySQL not yet reachable, retrying...", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			logger.Error("timeout waiting for MySQL", zap.Error(ctx.Err()))
			return nil, fmt.Errorf("could not connect to MySQL within 1m: %w", ctx.Err())
		case <-ticker.C:
			// retry
		}
	}

	// Register shutdown hook
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			logger.Info("closing MySQL connection")
			return db.Close()
		},
	})

	return db, nil
}

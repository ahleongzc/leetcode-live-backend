package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresDatabase(config *config.DatabaseConfig) (*sql.DB, error) {
	if config == nil {
		return nil, fmt.Errorf("missing database config: %w", common.ErrInternalServerError)
	}

	database, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to open database handler, %s: %w", err, common.ErrInternalServerError)
	}

	database.SetMaxOpenConns(config.MaxOpenConnections)
	database.SetMaxIdleConns(config.MaxIdleConnections)
	database.SetConnMaxIdleTime(config.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return database, nil
}

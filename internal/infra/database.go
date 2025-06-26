package infra

import (
	"context"
	"database/sql"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDatabase(config *config.DatabaseConfig) (*sql.DB, error) {
	database, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, err
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

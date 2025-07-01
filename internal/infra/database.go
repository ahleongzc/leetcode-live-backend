package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDatabase(
	config *config.DatabaseConfig,
) (*gorm.DB, error) {
	if config == nil {
		return nil, fmt.Errorf("missing database config: %w", common.ErrInternalServerError)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to open GORM database handler, %s: %w", err, common.ErrInternalServerError)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database from GORM DB, %s: %w", err, common.ErrInternalServerError)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

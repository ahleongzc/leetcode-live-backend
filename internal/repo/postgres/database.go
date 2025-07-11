package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/entity"

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
		Logger: logger.Default.LogMode(logger.Silent),
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

	err = migrateEntities(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateEntities(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Interview{},
		&entity.Question{},
		&entity.Session{},
		&entity.Transcript{},
		&entity.User{},
		&entity.Review{},
	)
	return err
}

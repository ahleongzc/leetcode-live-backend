package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type DatabaseConfig struct {
	DSN                string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        time.Duration
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	dsn := util.GetEnvOr(common.DB_DSN_KEY, "")

	maxOpenConnections := 25
	maxIdleConnections := 25
	maxIdleTime := 30 * time.Second

	maxOpenConnValue, err := strconv.Atoi(util.GetEnvOr(common.DB_MAX_OPEN_CONN_KEY, "25"))
	if nil == err {
		maxOpenConnections = maxOpenConnValue
	}

	maxIdleConnValue, err := strconv.Atoi(util.GetEnvOr(common.DB_MAX_IDLE_CONN_KEY, "25"))
	if nil == err {
		maxIdleConnections = maxIdleConnValue
	}

	maxIdleTimeSecondsValue, err := strconv.Atoi(util.GetEnvOr(common.DB_MAX_IDLE_TIME_SEC_KEY, "300"))
	if nil == err {
		maxIdleTime = time.Duration(maxIdleTimeSecondsValue) * time.Second
	}

	if dsn == "" {
		return nil, fmt.Errorf("database dsn cannot be empty: %w", common.ErrInternalServerError)
	}

	return &DatabaseConfig{
		DSN:                dsn,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
	}, nil
}

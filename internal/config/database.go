package config

import (
	"os"
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

func LoadDatabaseConfig() *DatabaseConfig {
	dsn := os.Getenv(common.DB_DSN_KEY)

	maxOpenConnections := 25
	maxIdleConnections := 25
	maxIdleTime := 30 * time.Second

	if util.IsProdEnv() {
		maxOpenConnValue, err := strconv.Atoi(os.Getenv(common.DB_MAX_OPEN_CONN_KEY))
		if nil == err {
			maxOpenConnections = maxOpenConnValue
		}

		maxIdleConnValue, err := strconv.Atoi(os.Getenv(common.DB_MAX_IDLE_CONN_KEY))
		if nil == err {
			maxIdleConnections = maxIdleConnValue
		}

		maxIdleTimeSecondsValue, err := strconv.Atoi(os.Getenv(common.DB_MAX_IDLE_TIME_SEC_KEY))
		if nil == err {
			maxIdleTime = time.Duration(maxIdleTimeSecondsValue) * time.Second
		}
	}

	return &DatabaseConfig{
		DSN:                dsn,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
	}
}

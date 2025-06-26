package config

import (
	"os"
	"strconv"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type ServerConfig struct {
	Address      string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadServerConfig() *ServerConfig {
	var address string
	var idleTimeout, readTimeout, writeTimeout time.Duration

	if util.IsDevEnv() {
		address = "localhost:8080"
		idleTimeout = time.Minute
		readTimeout = 10 * time.Second
		writeTimeout = 30 * time.Second
	}

	if util.IsProdEnv() {
		address = "localhost:" + os.Getenv(common.PORT_KEY)
		idleTimeout = toDurationSeconds(os.Getenv(common.IDLE_TIMEOUT_SEC_KEY))
		readTimeout = toDurationSeconds(os.Getenv(common.READ_TIMEOUT_SEC_KEY))
		writeTimeout = toDurationSeconds(os.Getenv(common.WRITE_TIMEOUT_SEC_KEY))
	}

	return &ServerConfig{
		Address:      address,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func toDurationSeconds(secondsString string) time.Duration {
	seconds, err := strconv.Atoi(secondsString)
	if err != nil {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

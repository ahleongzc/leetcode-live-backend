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
	address := "localhost:" + os.Getenv(common.PORT_KEY)
	idleTimeout := time.Minute
	readTimeout := 10 * time.Second
	writeTimeout := 30 * time.Second

	if util.IsProdEnv() {
		address = "0.0.0.0:" + os.Getenv(common.PORT_KEY)

		idleTimeoutSecondsValue, err := strconv.Atoi(os.Getenv(common.IDLE_TIMEOUT_SEC_KEY))
		if nil == err {
			idleTimeout = time.Duration(idleTimeoutSecondsValue) * time.Second
		}

		readTimeoutSecondsValue, err := strconv.Atoi(os.Getenv(common.READ_TIMEOUT_SEC_KEY))
		if nil == err {
			readTimeout = time.Duration(readTimeoutSecondsValue) * time.Second
		}

		writeTimeoutSecondsValue, err := strconv.Atoi(os.Getenv(common.WRITE_TIMEOUT_SEC_KEY))
		if nil == err {
			writeTimeout = time.Duration(writeTimeoutSecondsValue) * time.Second
		}
	}

	return &ServerConfig{
		Address:      address,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

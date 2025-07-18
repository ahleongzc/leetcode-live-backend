package config

import (
	"strconv"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HTTPServerConfig struct {
	Address      string
	Port         uint
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadHTTPServerConfig() *HTTPServerConfig {
	port := util.GetEnvUIntOr(common.HTTP_PORT_KEY, 8000)
	address := "0.0.0.0:" + util.GetEnvOr(common.HTTP_PORT_KEY, "8000")
	idleTimeout := time.Minute
	readTimeout := 10 * time.Second
	writeTimeout := 30 * time.Second

	if util.IsProdEnv() {
		idleTimeoutSecondsValue, err := strconv.Atoi(util.GetEnvOr(common.IDLE_TIMEOUT_SEC_KEY, "60"))
		if nil == err {
			idleTimeout = time.Duration(idleTimeoutSecondsValue) * time.Second
		}

		readTimeoutSecondsValue, err := strconv.Atoi(util.GetEnvOr(common.READ_TIMEOUT_SEC_KEY, "10"))
		if nil == err {
			readTimeout = time.Duration(readTimeoutSecondsValue) * time.Second
		}

		writeTimeoutSecondsValue, err := strconv.Atoi(util.GetEnvOr(common.WRITE_TIMEOUT_SEC_KEY, "30"))
		if nil == err {
			writeTimeout = time.Duration(writeTimeoutSecondsValue) * time.Second
		}
	}

	return &HTTPServerConfig{
		Port:         port,
		Address:      address,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

package util

import (
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

func IsDevEnv() bool {
	return GetEnvOr(common.ENVIRONMENT_KEY, common.DEV_ENVIRONMENT) == common.DEV_ENVIRONMENT
}

func IsProdEnv() bool {
	return GetEnvOr(common.ENVIRONMENT_KEY, common.DEV_ENVIRONMENT) == common.PROD_ENVIRONMENT
}

func GetEnvOr(envKey string, defaultValue string) string {
	value, ok := os.LookupEnv(envKey)
	if !ok {
		return defaultValue
	}
	if value == "" {
		return defaultValue
	}
	return value
}

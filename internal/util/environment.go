package util

import (
	"os"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

func IsDevEnv() bool {
	return os.Getenv(common.ENVIRONMENT_KEY) == common.DEV_ENVIRONMENT
}

func IsProdEnv() bool {
	return os.Getenv(common.ENVIRONMENT_KEY) == common.PROD_ENVIRONMENT
}

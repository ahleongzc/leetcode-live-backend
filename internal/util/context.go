package util

import (
	"context"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

func GetStartRequestTimestampMS(ctx context.Context) int64 {
	value := ctx.Value(common.REQUEST_TIMESTAMP_MS_CONTEXT_KEY)
	if value == nil {
		return 0
	}
	startRequestTimestampMS, ok := value.(int64)
	if !ok {
		return 0
	}
	return startRequestTimestampMS
}

func SetStartRequestTimestampMS(ctx context.Context) context.Context {
	return context.WithValue(ctx, common.REQUEST_TIMESTAMP_MS_CONTEXT_KEY, time.Now().UnixMilli())
}

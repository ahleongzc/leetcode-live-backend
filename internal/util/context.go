package util

import (
	"context"
	"fmt"
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

func SetSessionToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, common.SESSION_TOKEN_CONTEXT_KEY, token)
}

func GetSessionToken(ctx context.Context) (string, error) {
	value := ctx.Value(common.SESSION_TOKEN_CONTEXT_KEY)
	if value == nil {
		return "", fmt.Errorf("unable to get session token from context: %w", common.ErrInternalServerError)
	}
	sessionToken, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("unable to assert session token into string: %w", common.ErrInternalServerError)
	}
	return sessionToken, nil
}

func SetUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, common.USER_ID_CONTEXT_KEY, userID)
}

func GetUserID(ctx context.Context) (uint, error) {
	value := ctx.Value(common.USER_ID_CONTEXT_KEY)
	if value == nil {
		return 0, fmt.Errorf("unable to get user id from context: %w", common.ErrInternalServerError)
	}
	userID, ok := value.(uint)
	if !ok {
		return 0, fmt.Errorf("unable to assert user id into uint: %w", common.ErrInternalServerError)
	}
	return userID, nil
}

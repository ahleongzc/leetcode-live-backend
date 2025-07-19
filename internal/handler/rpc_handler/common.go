package rpchandler

import (
	"errors"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

func MapToRPCError(err error) error {
	switch {
	case errors.Is(err, common.ErrUnauthorized):
		return RPCErrUnauthorized
	case errors.Is(err, common.ErrForbidden):
		return RPCErrForbidden
	default:
		return RPCErrInternalServerError
	}
}

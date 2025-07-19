package rpchandler

import (
	"errors"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"

	"github.com/rs/zerolog/log"
)

func HandleErroResponseRPC(err error) error {
	switch {
	case errors.Is(err, common.ErrUnauthorized):
		return RPCErrUnauthorized
	case errors.Is(err, common.ErrForbidden):
		return RPCErrForbidden
	case errors.Is(err, common.ErrNotFound):
		return RPCErrNotFound
	default:
		log.Error().Err(err).Msg("error")
		return RPCErrInternalServerError
	}
}

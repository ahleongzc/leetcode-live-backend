package common

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrConflict            = errors.New("conflict")
	ErrRateLimited         = errors.New("rate limit exceeded")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrNormalClientClosure = errors.New("client closed connection normally")
)

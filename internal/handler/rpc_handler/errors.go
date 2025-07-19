package rpchandler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	RPCErrUnauthorized        = status.Error(codes.Unauthenticated, "unauthorized access")
	RPCErrForbidden           = status.Error(codes.PermissionDenied, "access forbidden")
	RPCErrNotFound            = status.Error(codes.NotFound, "not found")
	RPCErrInternalServerError = status.Error(codes.Internal, "internal server error")
)

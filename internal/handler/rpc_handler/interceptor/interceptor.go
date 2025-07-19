package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	rpchandler "github.com/ahleongzc/leetcode-live-backend/internal/handler/rpc_handler"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Interceptor struct {
	logger *zerolog.Logger
}

func NewInterceptor(
	logger *zerolog.Logger,
) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

// Unary
func (i *Interceptor) RecoverPanicUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			stackTrace := debug.Stack()
			i.logger.Error().
				Interface("panic", err).
				Bytes("stackTrace", stackTrace).
				Msg("panic recovered in rpc request")

			rpchandler.HandleErroResponseRPC(common.ErrInternalServerError)
		}
	}()

	return handler(ctx, req)
}

func (i *Interceptor) LoggerUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	i.logger.Info().
		Str("method", info.FullMethod).
		Msg("")

	return handler(ctx, req)
}

// Stream
func (i *Interceptor) LoggerStreamInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	i.logger.Info().
		Str("method", info.FullMethod).
		Bool("is_client_stream", info.IsClientStream).
		Bool("is_server_stream", info.IsServerStream).
		Msg("")

	return handler(srv, ss)
}

func (i *Interceptor) RecoverPanicStreamInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	defer func() {
		if err := recover(); err != nil {
			stackTrace := debug.Stack()
			i.logger.Error().
				Interface("panic", err).
				Bytes("stackTrace", stackTrace).
				Msg("panic recovered in rpc request")

			rpchandler.HandleErroResponseRPC(common.ErrInternalServerError)
		}
	}()

	return handler(srv, ss)
}

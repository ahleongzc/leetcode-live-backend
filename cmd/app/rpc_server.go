package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	rpchandler "github.com/ahleongzc/leetcode-live-backend/internal/handler/rpc_handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler/rpc_handler/interceptor"
	"github.com/ahleongzc/leetcode-live-backend/pb"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type RPCServer struct {
	srv             *grpc.Server
	logger          *zerolog.Logger
	rpcServerConfig *config.RPCServerConfig

	proxyHandler *rpchandler.ProxyHandler
	interceptor  *interceptor.Interceptor
}

func NewRPCServer(
	logger *zerolog.Logger,
	rpcServerConfig *config.RPCServerConfig,
	proxyHandler *rpchandler.ProxyHandler,
	interceptor *interceptor.Interceptor,
) *RPCServer {
	srv := grpc.NewServer(
		grpc.ConnectionTimeout(rpcServerConfig.ConnectionTimeout),
		grpc.MaxRecvMsgSize(int(rpcServerConfig.MaxRecvMsgSize)),
		grpc.MaxSendMsgSize(int(rpcServerConfig.MaxSendMsgSize)),
		grpc.ChainUnaryInterceptor(interceptor.RecoverPanicUnaryInterceptor, interceptor.LoggerUnaryInterceptor),
		grpc.ChainStreamInterceptor(interceptor.RecoverPanicStreamInterceptor, interceptor.LoggerStreamInterceptor),
	)

	return &RPCServer{
		srv:             srv,
		logger:          logger,
		proxyHandler:    proxyHandler,
		rpcServerConfig: rpcServerConfig,
		interceptor:     interceptor,
	}
}

func (rs *RPCServer) Serve(errChan chan error) *RPCServer {
	lis, err := net.Listen("tcp", rs.rpcServerConfig.Address)
	if err != nil {
		errChan <- err
		return nil
	}

	rs.registerHandlers()

	go func() {
		if err := rs.srv.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	rs.logger.Info().Msg(fmt.Sprintf("rpc server has started at port %d on %s", rs.rpcServerConfig.Port, time.Now().Format("2006-01-02 15:04:05")))
	return rs
}

func (rs *RPCServer) GracefullyTerminate(ctx context.Context) {
	if rs == nil || rs.srv == nil {
		return
	}
	rs.srv.GracefulStop()
	rs.logger.Info().Msg(fmt.Sprintf("rpc server has gracefully terminated at %s", time.Now().Format("2006-01-02 15:04:05")))
}

func (rs *RPCServer) registerHandlers() {
	if rs == nil || rs.srv == nil {
		return
	}
	pb.RegisterInterviewProxyServer(rs.srv, rs.proxyHandler)
}

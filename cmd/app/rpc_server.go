package app

import (
	"fmt"
	"net"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/proto"

	"google.golang.org/grpc"
)

type rpcServer struct {
	proto.UnimplementedInterviewProxyServer
}

func (a *Application) StartRPCServer(errChan chan error) *grpc.Server {
	rpcServerConfig := config.LoadRPCServerConfig()
	lis, err := net.Listen("tcp", rpcServerConfig.Address)
	if err != nil {
		errChan <- err
		return nil
	}

	srv := grpc.NewServer()
	proto.RegisterInterviewProxyServer(srv, &rpcServer{})

	go func() {
		if err := srv.Serve(lis); err != nil {
			errChan <- err
		}
	}()


	a.logger.Info().Msg(fmt.Sprintf("rpc server has started at %s", time.Now().Format("2006-01-02 15:04:05")))

	return srv
}

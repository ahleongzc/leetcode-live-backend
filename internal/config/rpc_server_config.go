package config

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type RPCServerConfig struct {
	Address string
}

func LoadRPCServerConfig() *RPCServerConfig {
	address := "0.0.0.0:" + util.GetEnvOr(common.RPC_PORT_KEY, "8100")

	return &RPCServerConfig{
		Address: address,
	}
}

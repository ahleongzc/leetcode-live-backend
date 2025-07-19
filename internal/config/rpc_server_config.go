package config

import (
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type RPCServerConfig struct {
	Address           string
	ConnectionTimeout time.Duration
	MaxRecvMsgSize    uint
	MaxSendMsgSize    uint
}

func LoadRPCServerConfig() *RPCServerConfig {
	address := "0.0.0.0:" + util.GetEnvOr(common.RPC_PORT_KEY, "8100")
	connectionTimeout := 10 * time.Second
	maxRecvMsgSize := uint(PAYLOAD_MAX_BYTES)
	maxSendMsgSize := uint(PAYLOAD_MAX_BYTES)

	return &RPCServerConfig{
		Address:           address,
		ConnectionTimeout: connectionTimeout,
		MaxRecvMsgSize:    maxRecvMsgSize,
		MaxSendMsgSize:    maxSendMsgSize,
	}
}

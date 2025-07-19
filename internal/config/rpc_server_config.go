package config

import (
	"strconv"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type RPCServerConfig struct {
	Port              uint
	Address           string
	ConnectionTimeout time.Duration
	MaxRecvMsgSize    uint
	MaxSendMsgSize    uint
}

func LoadRPCServerConfig() *RPCServerConfig {
	port := util.GetEnvUIntOr(common.RPC_PORT_KEY, 8100)
	address := "0.0.0.0:" + strconv.Itoa(int(port))
	connectionTimeout := 10 * time.Second
	maxRecvMsgSize := uint(PAYLOAD_MAX_BYTES)
	maxSendMsgSize := uint(PAYLOAD_MAX_BYTES)

	return &RPCServerConfig{
		Port:              port,
		Address:           address,
		ConnectionTimeout: connectionTimeout,
		MaxRecvMsgSize:    maxRecvMsgSize,
		MaxSendMsgSize:    maxSendMsgSize,
	}
}

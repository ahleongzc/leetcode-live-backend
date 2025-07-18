package rpchandler

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/proto"
)

type ProxyHandler struct{}

func (p *ProxyHandler) StartInterview(ctx context.Context, req *proto.StartInterviewRequest) (*proto.StartInterviewResponse, error) {
	return nil, nil
}

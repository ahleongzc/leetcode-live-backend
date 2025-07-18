package rpchandler

import (
	"context"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/proto"
)

func NewProxyHandler(
	interviewService service.InterviewService,
) *ProxyHandler {
	return &ProxyHandler{
		interviewService: interviewService,
	}
}

type ProxyHandler struct {
	proto.UnimplementedInterviewProxyServer
	interviewService service.InterviewService
}

func (p *ProxyHandler) VerifyCandidate(ctx context.Context, req *proto.VerifyCandidateRequest) (*proto.VerificationResponse, error) {
	token := req.GetToken()
	if token == "" {
		return nil, RPCErrUnauthorized
	}

	interview, err := p.interviewService.ConsumeTokenAndStartInterview(ctx, token)
	if err != nil {
		return nil, MapToRPCError(err)
	}

	resp := &proto.VerificationResponse{
		Interview: &proto.Interview{
			Id: uint64(interview.ID),
		},
	}

	return resp, nil
}

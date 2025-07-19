package rpchandler

import (
	"context"
	"io"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/pb"
)

func NewProxyHandler(
	interviewService service.InterviewService,
) *ProxyHandler {
	return &ProxyHandler{
		interviewService: interviewService,
	}
}

type ProxyHandler struct {
	pb.UnimplementedInterviewProxyServer
	interviewService service.InterviewService
}

func (p *ProxyHandler) ProcessIncomingMessage(stream pb.InterviewProxy_ProcessIncomingMessageServer) error {
	ctx := stream.Context()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return HandleErroResponseRPC(err)
		}

		res, err := p.interviewService.ProcessCandidateMessage(
			ctx,
			uint(in.GetInterviewId()),
			in.GetChunk(),
			in.GetCode(),
		)
		if err != nil {
			return HandleErroResponseRPC(err)
		}

		if !res.Exists() {
			continue
		}

		out := &pb.InterviewMessage{
			Source: pb.Source_SERVER,
			Url:    res.URL,
			End:    res.End,
		}

		if err := stream.Send(out); err != nil {
			return HandleErroResponseRPC(err)
		}
	}
}

func (p *ProxyHandler) VerifyCandidate(ctx context.Context, req *pb.VerifyCandidateRequest) (*pb.VerificationResponse, error) {
	token := req.GetToken()
	if token == "" {
		return nil, RPCErrUnauthorized
	}

	interview, err := p.interviewService.ConsumeTokenAndStartInterview(ctx, token)
	if err != nil {
		return nil, HandleErroResponseRPC(err)
	}

	resp := &pb.VerificationResponse{
		Interview: &pb.Interview{
			Id: uint64(interview.ID),
		},
	}

	return resp, nil
}

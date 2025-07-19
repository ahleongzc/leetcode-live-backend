package rpchandler

import (
	"context"
	"io"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
	"github.com/ahleongzc/leetcode-live-backend/pb"
)

func NewProxyHandler(
	authService service.AuthService,
	interviewService service.InterviewService,
) *ProxyHandler {
	return &ProxyHandler{
		authService:      authService,
		interviewService: interviewService,
	}
}

type ProxyHandler struct {
	pb.UnimplementedInterviewProxyServer
	authService      service.AuthService
	interviewService service.InterviewService
}

// TODO: See how to terminate this stream when the server is terminated as stream.Recv is a blocking operation
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
			Url:    util.ToPtr(res.URL),
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

	interviewID, err := p.authService.ValidateAndConsumeInterviewToken(ctx, token)
	if err != nil {
		return nil, HandleErroResponseRPC(err)
	}

	resp := &pb.VerificationResponse{
		InterviewId: util.ToPtr(uint64(interviewID)),
	}

	return resp, nil
}

func (p *ProxyHandler) JoinInterview(ctx context.Context, req *pb.JoinInterviewRequest) (*pb.JoinInterviewResponse, error) {
	interviewID := req.GetInterviewId()

	if err := p.interviewService.JoinInterview(ctx, uint(interviewID)); err != nil {
		return nil, HandleErroResponseRPC(err)
	}

	return &pb.JoinInterviewResponse{}, nil
}

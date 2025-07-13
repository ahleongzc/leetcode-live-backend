package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
)

type InterviewHandler struct {
	websocketConfig  *config.WebsocketConfig
	authService      service.AuthService
	interviewService service.InterviewService
	logger           *zerolog.Logger
}

func NewInterviewHandler(
	websocketConfig *config.WebsocketConfig,
	authService service.AuthService,
	interviewService service.InterviewService,
	logger *zerolog.Logger,
) *InterviewHandler {
	return &InterviewHandler{
		websocketConfig:  websocketConfig,
		authService:      authService,
		interviewService: interviewService,
		logger:           logger,
	}
}

func (i *InterviewHandler) EndOngoingInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	if err := i.interviewService.EndCandidateOngoingInterview(ctx, userID); err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	WriteJSONHTTP(w, nil, http.StatusOK, nil)
}

func (i *InterviewHandler) GetOngoingInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	ongoingInterview, err := i.interviewService.GetCandidateOngoingInterview(ctx, userID)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	payload := util.NewJSONPayload()
	payload.Add("data", util.JSONPayload{"interview": ongoingInterview})

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

func (i *InterviewHandler) AbandonUnfinishedInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	if err := i.interviewService.AbandonCandidateUnfinishedInterview(ctx, userId); err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	WriteJSONHTTP(w, nil, http.StatusOK, nil)
}

func (i *InterviewHandler) SetUpUnfinishedInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	token, err := i.interviewService.SetUpCandidateUnfinishedInterview(ctx, userID)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	header := http.Header{}
	header.Set(config.INTERVIEW_TOKEN_HEADER_KEY, token)

	WriteJSONHTTP(w, nil, http.StatusOK, header)
}

func (i *InterviewHandler) GetUnfinishedInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	unfinishedInterview, err := i.interviewService.GetCandidateUnfinishedInterview(ctx, userID)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	payload := util.NewJSONPayload()
	payload.Add("data", util.JSONPayload{"interview": unfinishedInterview})

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

func (i *InterviewHandler) GetInterviewHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	limit, offset := ParsePaginationParams(r)
	history, pagination, err := i.interviewService.GetHistory(ctx, userID, limit, offset)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	payload := util.NewJSONPayload()
	payload.Add("data", history)
	payload.Add("pagination", pagination)

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

// TODO: Don't allow user to set up new interview if there is too many abandoned interview since the last one
func (i *InterviewHandler) SetUpNewInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &struct {
		QuestionID  string `json:"question_id"`
		Description string `json:"description"`
	}{}

	err := ReadJSONHTTPReq(w, r, request)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	userID, err := util.GetUserID(ctx)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	token, err := i.interviewService.SetUpNewInterviewForCandidate(ctx, userID, request.QuestionID, request.Description)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	header := http.Header{}
	header.Set(config.INTERVIEW_TOKEN_HEADER_KEY, token)

	WriteJSONHTTP(w, nil, http.StatusOK, header)
}

func (i *InterviewHandler) JoinInterview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	token := r.URL.Query().Get("token")
	interview, err := i.interviewService.ConsumeTokenAndStartInterview(ctx, token)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	conn, err := websocket.Accept(w, r, i.websocketConfig.AcceptOptions)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "connection closed by server")

	respondChan := make(chan *model.WebSocketMessage)
	errChan := make(chan error)
	closeChan := make(chan struct{})

	go func() {
		defer close(respondChan)
		i.readPump(ctx, interview.ID, conn, respondChan, errChan)
	}()

	go func() {
		i.writePump(ctx, conn, respondChan, errChan, closeChan)
	}()

	select {
	case <-ctx.Done():
	case <-closeChan:
		conn.Close(websocket.StatusNormalClosure, "interview ended")
	case err := <-errChan:
		i.interviewService.PauseCandidateOngoingInterview(ctx, interview.UserID)
		HandleErrorResponeWebsocket(ctx, conn, err)
		cancel()
	}

	i.logger.Info().Msg(fmt.Sprintf("websocket connection closed for %s", r.RemoteAddr))
}

func (i *InterviewHandler) readPump(
	ctx context.Context,
	interviewID uint,
	conn *websocket.Conn,
	respondChan chan *model.WebSocketMessage,
	errChan chan error,
) {
	// Buffered channel for 20 incoming messages until it's processed downstream
	messageChan := make(chan *model.WebSocketMessage, 20)
	defer close(messageChan)

	go func() {
		for message := range messageChan {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := i.interviewService.ProcessIncomingMessage(ctx, interviewID, message)
				if err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done(): // To prevent writing to the error channel when the ctx is cancelled
					}
					continue
				}

				if response != nil {
					select {
					case respondChan <- response:
					case <-ctx.Done(): // To prevent writing to the respond channel when the ctx is cancelled
					}
				}
			}
		}
	}()

	for {
		_, bytes, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				err = common.ErrNormalClientClosure
			}
			errChan <- err
			return
		}

		websocketMessage := &model.WebSocketMessage{}
		err = ReadJSONBytes(bytes, websocketMessage)
		if err != nil {
			errChan <- err
			return
		}

		if !websocketMessage.ValidClientMessage() {
			errChan <- common.ErrBadRequest
			return
		}

		select {
		case messageChan <- websocketMessage:
		case <-ctx.Done():
			return
		}
	}
}

func (i *InterviewHandler) writePump(
	ctx context.Context,
	conn *websocket.Conn,
	respondChan <-chan *model.WebSocketMessage,
	errChan chan error,
	closeChan chan struct{},
) {
	for {
		select {
		case message, ok := <-respondChan:
			if !ok {
				return
			}
			payload := util.NewJSONPayload()
			payload.Add("from", message.From)
			payload.Add("url", message.URL)

			if err := WriteJSONWebsocket(ctx, conn, payload); err != nil {
				errChan <- err
				return
			}

			if message.CloseConn {
				closeChan <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}

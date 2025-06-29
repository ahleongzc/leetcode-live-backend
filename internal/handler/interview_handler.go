package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
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

type WebsocketMessageStruct struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (i *InterviewHandler) JoinInterview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	sessionID := r.URL.Query().Get("session_id")
	questionID := r.URL.Query().Get("question_id")

	valid, err := i.authService.ValidateSession(ctx, sessionID)
	if err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}
	if !valid {
		HandleErrorResponseHTTP(w, common.ErrUnauthorized)
		return
	}

	interview, err := i.interviewService.GetInterviewDetails(ctx, sessionID, questionID)
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

	respondChan := make(chan *model.InterviewMessage)
	errChan := make(chan error)

	go func() {
		defer close(respondChan)
		i.readPump(ctx, interview.ID, conn, respondChan, errChan)
	}()

	go func() {
		i.writePump(ctx, conn, respondChan, errChan)
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		HandleErrorResponeWebsocket(ctx, conn, err)
		cancel()
	}

	i.logger.Info().Msg(fmt.Sprintf("websocket connection closed for %s", r.RemoteAddr))
}

func (i *InterviewHandler) readPump(
	ctx context.Context,
	interviewID int,
	conn *websocket.Conn,
	respondChan chan *model.InterviewMessage,
	errChan chan error,
) {
	// Buffered channel for 20 incoming messages until it's processed downstream
	messageChan := make(chan *model.InterviewMessage, 20)
	defer close(messageChan)

	go func() {
		for message := range messageChan {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := i.interviewService.ProcessInterviewMessage(ctx, interviewID, message)
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

		websocketMessage := &WebsocketMessageStruct{}
		err = ReadJSONBytes(bytes, websocketMessage)
		if err != nil {
			errChan <- err
			return
		}

		message := &model.InterviewMessage{
			Type:    model.InterviewMessageType(websocketMessage.Type),
			Content: websocketMessage.Content,
		}
		select {
		case messageChan <- message:
		case <-ctx.Done():
			return
		}
	}
}

func (i *InterviewHandler) writePump(
	ctx context.Context,
	conn *websocket.Conn,
	respondChan <-chan *model.InterviewMessage,
	errChan chan error,
) {
	for {
		select {
		case message, ok := <-respondChan:
			if !ok {
				return
			}
			payload := util.NewJSONPayload()
			payload.Add("type", message.Type)
			payload.Add("content", message.Content)

			if err := WriteJSONWebsocket(ctx, conn, payload); err != nil {
				errChan <- err
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

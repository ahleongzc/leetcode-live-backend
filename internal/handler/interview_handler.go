package handler

import (
	"net/http"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"

	"github.com/coder/websocket"
)

type InterviewHandler struct {
	websocketConfig *config.WebsocketConfig
	authService     service.AuthService
}

func NewInterviewHandler(
	websocketConfig *config.WebsocketConfig,
	authService service.AuthService,
) *InterviewHandler {
	return &InterviewHandler{
		websocketConfig: websocketConfig,
		authService:     authService,
	}
}

type WebsocketMessageStruct struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (i *InterviewHandler) StartInterview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := r.Header.Get("session-id")

	valid, err := i.authService.ValidateSession(ctx, sessionID)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}

	if !valid {
		HandleErrorResponse(w, common.ErrUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, i.websocketConfig.AcceptOptions)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}
	defer conn.CloseNow()

	readChan := make(chan *model.InterviewMessage)
	writeChan := make(chan *model.InterviewMessage)

	defer close(readChan)
	defer close(writeChan)

	for {
		time.Sleep(25 * time.Millisecond)
	}
}

package model

import "github.com/ahleongzc/leetcode-live-backend/internal/util"

type Sender string

const (
	CLIENT Sender = "client"
	SERVER Sender = "server"
)

type WebSocketMessage struct {
	From      Sender  `json:"from"`
	Chunk     *string `json:"chunk"`
	Code      *string `json:"code"`
	URL       *string `json:"url"`
	CloseConn bool
}

func NewWebsocketMessage() *WebSocketMessage {
	return &WebSocketMessage{}
}

func NewServerWebsocketMessage() *WebSocketMessage {
	msg := NewWebsocketMessage().
		SetFrom(SERVER)
	return msg
}

func (w *WebSocketMessage) SetFrom(sender Sender) *WebSocketMessage {
	if w == nil {
		return nil
	}
	w.From = sender
	return w
}

func (w *WebSocketMessage) SetChunk(chunk string) *WebSocketMessage {
	if w == nil {
		return nil
	}
	w.Chunk = util.ToPtr(chunk)
	return w
}

func (w *WebSocketMessage) SetCode(code string) *WebSocketMessage {
	if w == nil {
		return nil
	}
	w.Code = util.ToPtr(code)
	return w
}

func (w *WebSocketMessage) SetURL(url string) *WebSocketMessage {
	if w == nil {
		return nil
	}
	w.URL = util.ToPtr(url)
	return w
}

func (w *WebSocketMessage) CloseConnection() *WebSocketMessage {
	if w == nil {
		return nil
	}
	w.CloseConn = true
	return w
}

func (w *WebSocketMessage) ValidClientMessage() bool {
	return w.From == CLIENT && w.URL == nil
}

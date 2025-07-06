package model

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

func (w *WebSocketMessage) ValidClientMessage() bool {
	return w.From == CLIENT && w.URL == nil
}

package config

import (
	"github.com/ahleongzc/leetcode-live-backend/internal/util"

	"github.com/coder/websocket"
)

type WebsocketConfig struct {
	AcceptOptions *websocket.AcceptOptions
}

func LoadWebsocketConfig() *WebsocketConfig {
	acceptOptions := &websocket.AcceptOptions{}

	if util.IsDevEnv() {
		acceptOptions.OriginPatterns = []string{"*"}
		acceptOptions.InsecureSkipVerify = true
	}

	if util.IsProdEnv() {
		acceptOptions.OriginPatterns = make([]string, 0)
		for origin := range PROD_TRUSTED_ORIGINS {
			acceptOptions.OriginPatterns = append(acceptOptions.OriginPatterns, origin)
		}
		acceptOptions.InsecureSkipVerify = false
	}

	return &WebsocketConfig{
		AcceptOptions: acceptOptions,
	}
}

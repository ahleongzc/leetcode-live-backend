package http

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
)

func NewHTTPCLient() *http.Client {
	return &http.Client{
		Timeout: config.HTTP_REQUEST_TIMEOUT,
	}
}

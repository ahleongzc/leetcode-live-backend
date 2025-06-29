package infra

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

func NewHTTPCLient() *http.Client {
	return &http.Client{
		Timeout: common.HTTP_REQUEST_TIMEOUT,
	}
}

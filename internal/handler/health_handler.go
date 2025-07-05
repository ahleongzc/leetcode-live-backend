package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hc *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := util.NewJSONPayload()
	payload.Add("health", "ok")

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

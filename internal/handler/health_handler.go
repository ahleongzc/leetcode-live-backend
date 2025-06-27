package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HealthHandler struct {
	tts  infra.TTS
	repo repo.FileRepo
}

func NewHealthHandler(
	tts infra.TTS,
	repo repo.FileRepo,
) *HealthHandler {
	return &HealthHandler{
		tts:  tts,
		repo: repo,
	}
}

func (hc *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := util.NewJSONPayload()
	payload.Add("health", "ok")

	WriteJSON(w, payload, http.StatusOK, nil)
}

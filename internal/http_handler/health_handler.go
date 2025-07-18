package httphandler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HealthHandler struct {
	transcriptManager service.TranscriptManager
}

func NewHealthHandler(
	transcriptManager service.TranscriptManager,
) *HealthHandler {
	return &HealthHandler{
		transcriptManager: transcriptManager,
	}
}

func (hc *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := util.NewJSONPayload()

	load := hc.transcriptManager.GetManagerInfo()

	payload.Add("health", "ok")
	payload.Add("concurrent users", load)

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

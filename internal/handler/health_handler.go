package handler

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HealthHandler struct {
	reviewScenario scenario.ReviewScenario
}

func NewHealthHandler(
	reviewScenario scenario.ReviewScenario,
) *HealthHandler {
	return &HealthHandler{
		reviewScenario: reviewScenario,
	}
}

func (hc *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := util.NewJSONPayload()
	payload.Add("health", "ok")

	if err := hc.reviewScenario.ReviewInterviewPerformance(r.Context(), 58); err != nil {
		HandleErrorResponseHTTP(w, err)
		return
	}

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

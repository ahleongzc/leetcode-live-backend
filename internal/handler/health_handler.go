package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type HealthHandler struct {
	producer infra.MessageQueueProducer
}

func NewHealthHandler(
	producer infra.MessageQueueProducer,
) *HealthHandler {
	return &HealthHandler{
		producer: producer,
	}
}

func (hc *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := util.NewJSONPayload()
	payload.Add("health", "ok")

	msg, _ := json.Marshal(payload)
	hc.producer.Push(r.Context(), msg, common.REVIEW_QUEUE)

	WriteJSONHTTP(w, payload, http.StatusOK, nil)
}

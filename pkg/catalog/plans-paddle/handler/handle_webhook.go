package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

func (h *Handler) PaddleWebhook(w http.ResponseWriter, r *http.Request) {
	if !h.IsEnabled {
		http.Error(w, "Paddle is not enabled", http.StatusServiceUnavailable)
		return
	}
	if h.PaddleWebhookSecret == "" {
		http.Error(w, "Paddle webhook secret is not configured", http.StatusServiceUnavailable)
		return
	}
	defer r.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
	if err != nil {
		http.Error(w, "unable to read Paddle webhook", http.StatusBadRequest)
		return
	}
	var event domain.PaddleWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil || event.EventType == "" {
		http.Error(w, "invalid Paddle webhook payload", http.StatusBadRequest)
		return
	}
	if err := h.PaddlePlanSrv.HandlePaddleWebhook(r.Context(), &event); err != nil {
		http.Error(w, "unable to process Paddle webhook", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

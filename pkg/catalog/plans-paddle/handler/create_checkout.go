package handler

import (
	"encoding/json"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type createCheckoutRequest struct {
	PlanID string `json:"plan_id"`
	ItemID string `json:"item_id"`
}

func (h *Handler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	if !h.IsEnabled {
		http.Error(w, "Paddle is not enabled", http.StatusServiceUnavailable)
		return
	}
	defer r.Body.Close()
	var req createCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PlanID == "" || req.ItemID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "plan_id and item_id are required"))
		return
	}
	checkout, err := h.PaddlePlanSrv.CreateCheckout(r.Context(), req.PlanID, req.ItemID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(checkout)
}

package handler

import (
	"encoding/json"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

// GetTransaction handles GET /paddle/transactions/{paddleTransactionID} when a client checks a Paddle transaction after checkout.
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	paddleTransactionID := r.PathValue("paddleTransactionID")
	if paddleTransactionID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing Paddle transaction id"))
		return
	}
	transaction, err := h.PaddlePlanSrv.GetTransaction(r.Context(), paddleTransactionID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transaction)
}

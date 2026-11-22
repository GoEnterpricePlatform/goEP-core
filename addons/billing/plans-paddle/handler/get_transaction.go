package handler

import (
	"encoding/json"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("paddleTransactionID")
	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing Paddle transaction id"))
		return
	}
	transaction, err := h.PaddleBillingSrv.GetTransaction(r.Context(), id)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transaction)
}

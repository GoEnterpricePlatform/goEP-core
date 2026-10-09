package handler

import (
	"encoding/json"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)


func (h *Handler) GetCheckout(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("transactionID")
	if transactionID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing transaction id"))
		return
	}
	checkout, err := h.PaddlePlanSrv.GetCheckout(r.Context(), transactionID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(checkout)
}
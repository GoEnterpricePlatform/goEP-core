package handler

import (
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing paddle plan id"))
		return
	}
	if err := h.PaddlePlanSrv.Delete(r.Context(), id); err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

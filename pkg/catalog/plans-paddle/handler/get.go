package handler

import (
	"encoding/json"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Get handles GET /paddle/plans/{id} when a client needs the public catalog details for one Paddle plan.
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing paddle plan id"))
		return
	}
	post, err := h.PaddlePlanSrv.Get(r.Context(), id)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

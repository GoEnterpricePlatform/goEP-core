package handler

import (
	"context"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) DeleteVarOption(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	variationID := r.PathValue("variationId")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing varOption id"))
		return
	}

	if variationID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing variationId id"))
		return
	}

	if err := h.VariationSrv.DeleteVarOption(context.Background(), id, variationID); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
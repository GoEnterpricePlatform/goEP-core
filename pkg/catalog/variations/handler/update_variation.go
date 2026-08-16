package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) UpdateVariation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing category id"))
		return
	}

	var req core.UpdateVariationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	variation := &domain.Variation{
		Name: req.Name,
		ID:   id,
	}

	if err := h.VariationSrv.UpdateVariation(context.Background(), variation); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(variation)

}

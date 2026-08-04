package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) CreateVarOption(w http.ResponseWriter, r *http.Request) {
	variationID := r.PathValue("variationId")

	if variationID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing variation id"))
		return
	}

	var req core.CreateVarOptionReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	varOption := &domain.VarOption{
		Label:       req.Label,
		Value:       req.Value,
		VariationID: variationID,
	}

	if err := h.VariationSrv.CreateVarOption(context.Background(), varOption); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(varOption)
}

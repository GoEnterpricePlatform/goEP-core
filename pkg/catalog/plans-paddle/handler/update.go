package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing paddle plan id"))
		return
	}
	defer r.Body.Close()
	var req core.UpdatePaddlePlanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}
	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}
	paddlePlan := &domain.PaddlePlan{Name: req.Name, Description: req.Description, PaddleProductID: req.PaddleProductID, Order: req.Order, Items: make([]*domain.PaddlePlanItem, 0, len(req.Items))}
	for _, item := range req.Items {
		paddlePlan.Items = append(paddlePlan.Items, &domain.PaddlePlanItem{ID: item.ID, Features: item.Features, VarOptionIDs: item.VarOptionIDs, PaddlePriceID: item.PaddlePriceID})
	}
	if err := h.PaddlePlanSrv.Update(r.Context(), id, paddlePlan); err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paddlePlan)
}

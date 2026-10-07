package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Create a product, with or without variations, and establish your relationships with Paddle.
// The returned product may be incomplete;
// uses your ID to get full plan information using /plans/{id}.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req core.CreatePaddlePlanReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	paddlePlan := &domain.PaddlePlan{
		Name:            req.Name,
		Description:     req.Description,
		PaddleProductID: req.PaddleProductID,
		Order:           req.Order,
		Items:           make([]*domain.PaddlePlanItem, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		paddlePlan.Items = append(
			paddlePlan.Items,
			&domain.PaddlePlanItem{
				Features:      item.Features,
				VarOptionIDs:  item.VarOptionIDs,
				PaddlePriceID: item.PaddlePriceID,
			},
		)
	}

	if err := h.PaddlePlanSrv.Create(r.Context(), paddlePlan); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paddlePlan)
}

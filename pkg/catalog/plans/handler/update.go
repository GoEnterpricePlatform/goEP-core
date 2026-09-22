package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(
			w,
			sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing plan id"),
		)
		return
	}

	var req core.CreateOrUpdatePlanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(
			w,
			sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"),
		)
		return
	}

	defer r.Body.Close()

	if err := req.ValidateUpdate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	plan := &domain.Plan{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Items:       make([]*domain.PlanItem, 0, len(req.PlanItems)),
	}

	for _, item := range req.PlanItems {
		plan.Items = append(plan.Items, &domain.PlanItem{
			ID:           item.ID,
			Features:     item.Features,
			VarOptionIDs: item.VarOptionIDs,
		})
	}

	if err := h.PlanSrv.Update(context.Background(), id, plan); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(plan)
}

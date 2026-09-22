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

func (h Handler) Patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(
			sharedD.ErrCodeInvalidParams,
			"missing plan id",
		),
		)
		return
	}

	var req core.PatchPlanReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(
			sharedD.ErrCodeInvalidParams,
			"invalid request body",
		),
		)
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	plan := &domain.Plan{}

	if req.Name != nil {
		plan.Name = *req.Name
	}

	plan.Description = req.Description

	if req.PlanItems != nil {
		plan.Items = make([]*domain.PlanItem, 0, len(*req.PlanItems))

		for _, itemReq := range *req.PlanItems {
			item := &domain.PlanItem{}

			item.ID = itemReq.ID

			if itemReq.Features != nil {
				item.Features = *itemReq.Features
			}

			if itemReq.VarOptionIDs != nil {
				item.VarOptionIDs = *itemReq.VarOptionIDs
			}

			plan.Items = append(plan.Items, item)
		}
	}

	updated, err := h.PlanSrv.Patch(context.Background(), id, plan)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

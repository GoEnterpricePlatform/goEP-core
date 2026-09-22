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

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req core.CreateOrUpdatePlanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	plan := &domain.Plan {
		Name:        req.Name,
		Description: req.Description,
	}

	// Plan with variants
	if len(req.PlanItems) > 1 {
		for _, pItem := range req.PlanItems {
			pItem := &domain.PlanItem {
				Features:     pItem.Features,
				VarOptionIDs: pItem.VarOptionIDs,
			}
			plan.Items = append(plan.Items, pItem)
		}
	} else {
		// plan without variants -> a single item by default
		pItem := &domain.PlanItem{
			Features:     req.PlanItems[0].Features,
			VarOptionIDs: req.PlanItems[0].VarOptionIDs,
		}
		plan.Items = append(plan.Items, pItem)
	}

	if err := h.PlanSrv.Create(context.Background(), plan); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plan)
}

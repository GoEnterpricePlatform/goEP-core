package handler

import (
	"net/http"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"github.com/starfederation/datastar-go/datastar"
	frontCore "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
)

func (h *Handler) CreatePaddlePlan(w http.ResponseWriter, r *http.Request) {
	h.savePaddlePlan(w, r, "")
}

func (h *Handler) UpdatePaddlePlan(w http.ResponseWriter, r *http.Request) {
	h.savePaddlePlan(w, r, r.PathValue("id"))
}

func (h *Handler) savePaddlePlan(w http.ResponseWriter, r *http.Request, id string) {
	var signals paddlePlanSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	f := signals.PlanForm
	name := strings.TrimSpace(f.Name)
	productID := strings.TrimSpace(f.PaddleProductID)
	if name == "" || productID == "" || len(f.Items) == 0 {
		h.patchError(w, r, "Name, Paddle Product ID and at least one item are required")
		return
	}
	p := &domain.PaddlePlan{Name: name, PaddleProductID: productID, Order: f.Order, Items: make([]*domain.PaddlePlanItem, 0, len(f.Items))}
	if description := strings.TrimSpace(f.Description); description != "" {
		p.Description = &description
	}
	for _, item := range f.Items {
		features := cleanFeatures(item.Features)
		if strings.TrimSpace(item.PaddlePriceID) == "" {
			h.patchError(w, r, "Every combination needs a Paddle Price ID")
			return
		}
		if len(features) == 0 {
			h.patchError(w, r, "Add at least one feature to every combination")
			return
		}
		p.Items = append(p.Items, &domain.PaddlePlanItem{ID: item.ID, Features: features, VarOptionIDs: item.VarOptionIDs, PaddlePriceID: strings.TrimSpace(item.PaddlePriceID), Status: planDomain.PlanStatus(item.Status)})
	}
	var err error
	if id == "" {
		err = h.PaddlePlanSrv.Create(r.Context(), p)
	} else {
		err = h.PaddlePlanSrv.Update(r.Context(), id, p)
	}
	if err != nil {
		h.patchError(w, r, frontCore.UiErrorResp(err))
		return
	}
	sse := datastar.NewSSE(w, r)
	_ = sse.Redirect("/v1/catalog/paddle-plans")
}

func cleanFeatures(rawFeatures []string) []string {
	result := make([]string, 0, len(rawFeatures))
	for _, raw := range rawFeatures {
		if value := strings.TrimSpace(raw); value != "" {
			result = append(result, value)
		}
	}
	return result
}

package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/components"
	frontCore "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
	"github.com/starfederation/datastar-go/datastar"
)

func (h *Handler) DeletePaddlePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.patchError(w, r, "Plan ID is required")
		return
	}
	if err := h.PaddlePlanSrv.Delete(r.Context(), id); err != nil {
		h.patchError(w, r, frontCore.UiErrorResp(err))
		return
	}
	plans, _, _, err := h.PaddlePlanSrv.GetAll(r.Context(), 100, 1)
	if err != nil {
		h.patchError(w, r, frontCore.UiErrorResp(err))
		return
	}
	sse := datastar.NewSSE(w, r)
	_ = sse.MarshalAndPatchSignals(map[string]any{
		"error":      "",
		"is_loading": false,
		"modal":      map[string]any{"domain": "", "action": ""},
		"plan_form":  map[string]any{"id": ""},
	})
	_ = sse.PatchElementTempl(components.PaddlePlanList(plans), datastar.WithSelector("#paddle-plan-list"), datastar.WithModeOuter())
}

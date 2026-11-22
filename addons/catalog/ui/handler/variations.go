package handler

import (
	"net/http"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/components"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/core"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
	"github.com/starfederation/datastar-go/datastar"
)

type variationFormSignals struct {
	VariationForm struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"variation_form"`
}

func (h *Handler) CreateVariation(w http.ResponseWriter, r *http.Request) {
	var signals variationFormSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	req := core.CreateVariationReq{Name: strings.TrimSpace(signals.VariationForm.Name)}
	if err := req.Validate(); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	if err := h.VariationSrv.CreateVariation(r.Context(), &domain.Variation{Name: req.Name}); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"variation_form": map[string]any{"id": "", "name": ""},
	})
}

func (h *Handler) UpdateVariation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var signals variationFormSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	req := core.UpdateVariationReq{Name: strings.TrimSpace(signals.VariationForm.Name)}
	if id == "" {
		h.patchError(w, r, "Variation ID is required")
		return
	}
	if err := req.Validate(); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	if err := h.VariationSrv.UpdateVariation(r.Context(), &domain.Variation{ID: id, Name: req.Name}); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"variation_form": map[string]any{"id": "", "name": ""},
	})
}

func (h *Handler) DeleteVariation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.patchError(w, r, "Variation ID is required")
		return
	}
	if err := h.VariationSrv.DeleteVariation(r.Context(), id); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"variation_form": map[string]any{"id": "", "name": ""},
	})
}

func (h *Handler) refreshVariations(w http.ResponseWriter, r *http.Request, resetSignals map[string]any) {
	variations, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
	if err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	sse := datastar.NewSSE(w, r)
	signals := map[string]any{
		"error":      "",
		"is_loading": false,
		"modal":      map[string]any{"domain": "", "action": ""},
	}
	for key, value := range resetSignals {
		signals[key] = value
	}
	if err := sse.MarshalAndPatchSignals(signals); err != nil {
		return
	}
	_ = sse.PatchElementTempl(
		components.VariationList(variations),
		datastar.WithSelector("#variation-list"),
		datastar.WithModeOuter(),
	)
}

func (h *Handler) patchError(w http.ResponseWriter, r *http.Request, message string) {
	sse := datastar.NewSSE(w, r)
	_ = sse.MarshalAndPatchSignals(map[string]any{
		"error":      message,
		"is_loading": false,
	})
}

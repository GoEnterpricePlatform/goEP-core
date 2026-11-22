package handler

import (
	"net/http"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/core"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
	"github.com/starfederation/datastar-go/datastar"
)

type optionFormSignals struct {
	OptionForm struct {
		ID          string `json:"id"`
		VariationID string `json:"variation_id"`
		Label       string `json:"label"`
		Value       string `json:"value"`
	} `json:"option_form"`
}

func (h *Handler) CreateOption(w http.ResponseWriter, r *http.Request) {
	var signals optionFormSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	variationID := r.PathValue("variationId")
	req := core.CreateVarOptionReq{Label: strings.TrimSpace(signals.OptionForm.Label), Value: optionalValue(signals.OptionForm.Value)}
	if variationID == "" {
		h.patchError(w, r, "Variation ID is required")
		return
	}
	if err := req.Validate(); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	if err := h.VariationSrv.CreateVarOption(r.Context(), &domain.VarOption{
		VariationID: variationID,
		Label:       req.Label,
		Value:       req.Value,
	}); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"option_form": map[string]any{"id": "", "variation_id": "", "variation_name": "", "label": "", "value": ""},
	})
}

func (h *Handler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	var signals optionFormSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	variationID := r.PathValue("variationId")
	id := r.PathValue("id")
	req := core.UpdateVarOptionReq{Label: strings.TrimSpace(signals.OptionForm.Label), Value: optionalValue(signals.OptionForm.Value)}
	if variationID == "" || id == "" {
		h.patchError(w, r, "Variation and option IDs are required")
		return
	}
	if err := req.Validate(); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	if err := h.VariationSrv.UpdateVarOption(r.Context(), &domain.VarOption{
		ID:          id,
		VariationID: variationID,
		Label:       req.Label,
		Value:       req.Value,
	}); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"option_form": map[string]any{"id": "", "variation_id": "", "variation_name": "", "label": "", "value": ""},
	})
}

func (h *Handler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	variationID := r.PathValue("variationId")
	id := r.PathValue("id")
	if variationID == "" || id == "" {
		h.patchError(w, r, "Variation and option IDs are required")
		return
	}
	if err := h.VariationSrv.DeleteVarOption(r.Context(), id, variationID); err != nil {
		h.patchError(w, r, sharedC.UiErrorResp(err))
		return
	}
	h.refreshVariations(w, r, map[string]any{
		"option_form": map[string]any{"id": "", "variation_id": "", "variation_name": "", "label": "", "value": ""},
	})
}

func optionalValue(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

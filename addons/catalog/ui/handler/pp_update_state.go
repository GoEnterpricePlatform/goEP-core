package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/components"
	variationDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/domain"
	frontCore "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
	"github.com/starfederation/datastar-go/datastar"
)

// Triggered when changing the product type or checking/unchecking an option.
// GeneratePaddlePlanItems reads the signals, recalculates the joins,
// and updates the Items table.
func (h *Handler) GeneratePaddlePlanItems(w http.ResponseWriter, r *http.Request) {
	var signals paddlePlanSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	variations, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
	if err != nil {
		h.patchError(w, r, frontCore.UiErrorResp(err))
		return
	}
	selected := make(map[string]bool, len(signals.PlanForm.SelectedOptionIDs))
	for _, id := range signals.PlanForm.SelectedOptionIDs {
		selected[id] = true
	}
	groups := make([][]*variationDomain.VarOption, 0)
	labels := make(map[string]string)
	for _, variation := range variations {
		var group []*variationDomain.VarOption
		for _, option := range variation.Options {
			if selected[option.ID] {
				group = append(group, option)
				labels[option.ID] = variation.Name + ": " + option.Label
			}
		}
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}
	combinations := [][]string{{}}
	for _, group := range groups {
		next := make([][]string, 0, len(combinations)*len(group))
		for _, combination := range combinations {
			for _, option := range group {
				next = append(next, append(append([]string{}, combination...), option.ID))
			}
		}
		combinations = next
	}
	if signals.PlanForm.ProductType != "variants" || len(combinations) <= 1 {
		combinations = [][]string{{}}
	}
	previous := make(map[string]paddlePlanItemSignals)
	for _, item := range signals.PlanForm.Items {
		previous[optionKey(item.VarOptionIDs)] = item
	}
	items := make([]paddlePlanItemSignals, 0, len(combinations))
	views := make([]components.PlanItemForm, 0, len(combinations))
	for _, combination := range combinations {
		key := optionKey(combination)
		item, ok := previous[key]
		if !ok {
			item = paddlePlanItemSignals{Features: []string{""}, Status: "active"}
		}
		if len(item.Features) == 0 {
			item.Features = []string{""}
		}
		item.VarOptionIDs = combination
		items = append(items, item)
		parts := make([]string, 0, len(combination))
		for _, optionID := range combination {
			parts = append(parts, labels[optionID])
		}
		label := "Single plan"
		if len(parts) > 0 {
			label = strings.Join(parts, " / ")
		}
		views = append(views, components.PlanItemForm{Combination: label, Features: item.Features})
	}
	sse := datastar.NewSSE(w, r)
	_ = sse.MarshalAndPatchSignals(map[string]any{"error": "", "is_loading": false, "plan_form": map[string]any{"items": items}})
	_ = sse.PatchElementTempl(components.PaddlePlanItems(views), datastar.WithSelector("#paddle-plan-items"), datastar.WithModeOuter())
}

// Add a feature field
func (h *Handler) AddPaddlePlanFeature(w http.ResponseWriter, r *http.Request) {
	h.mutatePaddlePlanFeature(w, r, true)
}

// Remove a feature field
func (h *Handler) RemovePaddlePlanFeature(w http.ResponseWriter, r *http.Request) {
	h.mutatePaddlePlanFeature(w, r, false)
}

func (h *Handler) mutatePaddlePlanFeature(w http.ResponseWriter, r *http.Request, add bool) {
	var signals paddlePlanSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		h.patchError(w, r, "Invalid request")
		return
	}
	itemIndex, err := strconv.Atoi(r.PathValue("itemIndex"))
	if err != nil || itemIndex < 0 || itemIndex >= len(signals.PlanForm.Items) {
		h.patchError(w, r, "Plan item is required")
		return
	}
	item := &signals.PlanForm.Items[itemIndex]
	if add {
		item.Features = append(item.Features, "")
	} else {
		featureIndex, parseErr := strconv.Atoi(r.PathValue("featureIndex"))
		if parseErr != nil || featureIndex < 0 || featureIndex >= len(item.Features) {
			h.patchError(w, r, "Feature is required")
			return
		}
		item.Features = append(item.Features[:featureIndex], item.Features[featureIndex+1:]...)
		if len(item.Features) == 0 {
			item.Features = []string{""}
		}
	}
	variations, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
	if err != nil {
		h.patchError(w, r, frontCore.UiErrorResp(err))
		return
	}
	views := make([]components.PlanItemForm, 0, len(signals.PlanForm.Items))
	for _, planItem := range signals.PlanForm.Items {
		parts := make([]string, 0, len(planItem.VarOptionIDs))
		for _, variation := range variations {
			for _, option := range variation.Options {
				for _, selectedID := range planItem.VarOptionIDs {
					if option.ID == selectedID {
						parts = append(parts, variation.Name+": "+option.Label)
					}
				}
			}
		}
		label := "Single plan"
		if len(parts) > 0 {
			label = strings.Join(parts, " / ")
		}
		views = append(views, components.PlanItemForm{Combination: label, Features: planItem.Features})
	}
	sse := datastar.NewSSE(w, r)
	_ = sse.MarshalAndPatchSignals(map[string]any{"error": "", "is_loading": false, "plan_form": map[string]any{"items": signals.PlanForm.Items}})
	_ = sse.PatchElementTempl(components.PaddlePlanItems(views), datastar.WithSelector("#paddle-plan-items"), datastar.WithModeOuter())
}

func optionKey(ids []string) string {
	sorted := append([]string{}, ids...)
	sort.Strings(sorted)
	return strings.Join(sorted, ":")
}
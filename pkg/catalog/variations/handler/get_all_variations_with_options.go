package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
)

func (h Handler) GetAllVariationsWithOptions(w http.ResponseWriter, r *http.Request) {
	categories, err := h.VariationSrv.GetAllVariationsWithOptions(r.Context())
	if err != nil {
		core.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
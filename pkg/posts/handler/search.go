package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Search handles the HTTP request to search posts using a query string
// and an optional limit parameter.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	query := queryParams.Get("query")
	limitStr := queryParams.Get("limit")

	if query == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing query parameter"))
		return
	}

	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}

	posts, err := h.PostSrv.Search(r.Context(), query, limit)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

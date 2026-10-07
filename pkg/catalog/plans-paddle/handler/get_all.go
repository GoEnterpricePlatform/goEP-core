package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
)

type paddlePlansPage struct {
	Count    int64                `json:"count"`
	Pages    int64                `json:"pages"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Plans    []*domain.PaddlePlan `json:"plans"`
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	plans, count, totalPages, err := h.PaddlePlanSrv.GetAll(
		r.Context(),
		int64(limit),
		int64(page),
	)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	// Determine request scheme.
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host
	basePath := "/v1" + r.URL.Path

	// Build pagination links.
	var nextURL *string
	if int64(page) < totalPages {
		url := fmt.Sprintf(
			"%s://%s%s?page=%d&limit=%d",
			scheme,
			host,
			basePath,
			page+1,
			limit,
		)
		nextURL = &url
	}

	var prevURL *string
	if page > 1 {
		url := fmt.Sprintf(
			"%s://%s%s?page=%d&limit=%d",
			scheme,
			host,
			basePath,
			page-1,
			limit,
		)
		prevURL = &url
	}

	resp := &paddlePlansPage{
		Count:    count,
		Pages:    totalPages,
		Next:     nextURL,
		Previous: prevURL,
		Plans:    plans,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}

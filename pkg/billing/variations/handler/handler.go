package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/port"
)

type Handler struct {
	VariationSrv port.VariationSrv
}

func NewVariationHandler(mux *http.ServeMux, variationSrv port.VariationSrv) *Handler {
	h := &Handler{
		VariationSrv: variationSrv,
	}

	mux.HandleFunc("POST /variations", h.CreateVariation)
	mux.HandleFunc("GET /variations/options", h.GetAllVariationsWithOptions)
	mux.HandleFunc("PUT /variations/{id}", h.UpdateVariation)
	mux.HandleFunc("DELETE /variations/{id}", h.DeleteVariation)

	mux.HandleFunc("POST /variations/{variationId}/options", h.CreateVarOption)
	mux.HandleFunc("PUT /variations/{variationId}/options/{id}", h.UpdateVarOption)
	mux.HandleFunc("DELETE /variations/{variationId}/options/{id}", h.DeleteVarOption)

	return h
}

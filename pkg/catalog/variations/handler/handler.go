package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	VariationSrv port.VariationSrv
	AuthApiMdw   *middlewares.AuthMiddleware
}

func NewVariationHandler(
	mux *http.ServeMux,
	variationSrv port.VariationSrv,
	authApiMdw *middlewares.AuthMiddleware,
) *Handler {
	h := &Handler{
		VariationSrv: variationSrv,
		AuthApiMdw: authApiMdw,
	}

	mux.Handle("POST /variations", h.AuthApiMdw.AccessTokenMdw(h.CreateVariation))
	mux.HandleFunc("GET /variations/options", h.GetAllVariationsWithOptions)
	mux.Handle("PUT /variations/{id}", h.AuthApiMdw.AccessTokenMdw(h.UpdateVariation))
	mux.Handle("DELETE /variations/{id}", h.AuthApiMdw.AccessTokenMdw(h.DeleteVariation))

	mux.Handle("POST /variations/{variationId}/options", h.AuthApiMdw.AccessTokenMdw(h.CreateVarOption))
	mux.Handle("PUT /variations/{variationId}/options/{id}", h.AuthApiMdw.AccessTokenMdw(h.UpdateVarOption))
	mux.Handle("DELETE /variations/{variationId}/options/{id}", h.AuthApiMdw.AccessTokenMdw(h.DeleteVarOption))

	return h
}

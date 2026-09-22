package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	PlanSrv    port.PlanSrv
	AuthApiMdw *middlewares.AuthMiddleware
}

func NewPlanHandler(muxV1 *http.ServeMux, planSrv port.PlanSrv, authApiMdw *middlewares.AuthMiddleware) *Handler {
	h := &Handler{
		PlanSrv: planSrv,
		AuthApiMdw: authApiMdw,
	}

	muxV1.HandleFunc("GET /plans", h.GetAll)
	muxV1.Handle("POST /plans", h.AuthApiMdw.AccessTokenMdw(h.Create))
	muxV1.HandleFunc("GET /plans/{id}", h.Get)
	muxV1.Handle("PUT /plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.Update))
	muxV1.Handle("PATCH /plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.Patch))

	// muxV1.Handle("POST /posts", h.AuthApiMdw.AccessTokenMdw(h.Create))
	// muxV1.HandleFunc("GET /posts/{id}", h.Get)
	// muxV1.HandleFunc("GET /posts", h.GetAll)
	// muxV1.Handle("PUT /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Update))
	// muxV1.Handle("PATCH /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Patch))
	// muxV1.Handle("DELETE /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Delete))
	// muxV1.HandleFunc("GET /posts/search", h.Search)

	return h
}

package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	PaddlePlanSrv port.PaddlePlanSrv
	AuthApiMdw    *middlewares.AuthMiddleware
}

func NewPlanPaddleHandler(muxV1 *http.ServeMux, paddlePlanSrv port.PaddlePlanSrv, authApiMdw *middlewares.AuthMiddleware) *Handler {

	h := &Handler{
		PaddlePlanSrv: paddlePlanSrv,
		AuthApiMdw:    authApiMdw,
	}

	muxV1.Handle("POST /paddle/plans", h.AuthApiMdw.AccessTokenMdw(h.Create))
	muxV1.HandleFunc("GET /paddle/plans/{id}", h.Get)
	muxV1.HandleFunc("GET /paddle/plans", h.GetAll)
	muxV1.Handle("PUT /paddle/plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.Update))
	muxV1.Handle("DELETE /paddle/plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.Delete))

	return h
}

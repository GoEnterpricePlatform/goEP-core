package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	PProviderSrv port.PaymentProviderSrv
	AuthApiMdw   *middlewares.AuthMiddleware
}

func NewPaymentProviderApiHandler(
	muxV1 *http.ServeMux,
	pProviderSrv port.PaymentProviderSrv,
	authApiMdw *middlewares.AuthMiddleware,
) *Handler {
	h := &Handler{
		PProviderSrv: pProviderSrv,
		AuthApiMdw:   authApiMdw,
	}

	muxV1.Handle("POST /payments/provider", h.AuthApiMdw.AccessTokenMdw(h.Create))

	return h
}

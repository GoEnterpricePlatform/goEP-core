package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/payment-providers/port"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
)

type Handler struct {
	PProviderSrv port.PaymentProviderSrv
	AuthApiMdw   middlewarePort.Authenticator
}

func NewPaymentProviderApiHandler(
	muxV1 *http.ServeMux,
	pProviderSrv port.PaymentProviderSrv,
	authApiMdw middlewarePort.Authenticator,
) *Handler {
	h := &Handler{
		PProviderSrv: pProviderSrv,
		AuthApiMdw:   authApiMdw,
	}

	muxV1.Handle("POST /payments/provider", h.AuthApiMdw.AccessTokenMdw(h.Create))

	return h
}

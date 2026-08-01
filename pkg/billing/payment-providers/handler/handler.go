package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/port"
)

type Handler struct {
	PProviderSrv port.PaymentProviderSrv
}

func NewPaymentProviderApiHandler(muxV1 *http.ServeMux, pProviderSrv port.PaymentProviderSrv) *Handler {
	h := &Handler{
		PProviderSrv: pProviderSrv,
	}

	muxV1.HandleFunc("POST /payments/provider", h.Create)

	return h
}

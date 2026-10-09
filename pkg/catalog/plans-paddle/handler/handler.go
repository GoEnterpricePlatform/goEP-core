package handler

import (
	"net/http"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type Handler struct {
	PaddlePlanSrv       port.PaddlePlanSrv
	AuthApiMdw          *middlewares.AuthMiddleware
	IsEnabled           bool
	PaddleWebhookSecret string
}

func NewPlanPaddleHandler(muxV1 *http.ServeMux, paddlePlanSrv port.PaddlePlanSrv, authApiMdw *middlewares.AuthMiddleware, webhookSecret string, isEnabled bool) *Handler {

	h := &Handler{
		PaddlePlanSrv:       paddlePlanSrv,
		AuthApiMdw:          authApiMdw,
		IsEnabled:           isEnabled,
		PaddleWebhookSecret: webhookSecret,
	}

	muxV1.Handle("POST /paddle/plans", h.AuthApiMdw.AccessTokenMdw(h.requirePaddleEnabled(h.Create)))
	muxV1.HandleFunc("GET /paddle/plans/{id}", h.requirePaddleEnabled(h.Get))
	muxV1.HandleFunc("GET /paddle/plans", h.requirePaddleEnabled(h.GetAll))
	muxV1.Handle("PUT /paddle/plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.requirePaddleEnabled(h.Update)))
	muxV1.Handle("DELETE /paddle/plans/{id}", h.AuthApiMdw.AccessTokenMdw(h.requirePaddleEnabled(h.Delete)))
	muxV1.HandleFunc("POST /paddle/checkout", h.requirePaddleEnabled(h.CreateCheckout))
	muxV1.HandleFunc("GET /paddle/checkout/{transactionID}", h.requirePaddleEnabled(h.GetCheckout))
	var webhookHandler http.Handler = http.HandlerFunc(h.PaddleWebhook)
	if isEnabled && webhookSecret != "" {
		webhookHandler = paddleSDK.NewWebhookVerifier(webhookSecret, paddleSDK.VerifierWithTimestampTolerance(5*time.Second)).Middleware(webhookHandler)
	}
	muxV1.Handle("POST /paddle/webhook", webhookHandler)

	return h
}

func (h *Handler) requirePaddleEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.IsEnabled {
			http.Error(w, "Paddle is not enabled", http.StatusServiceUnavailable)
			return
		}
		next(w, r)
	}
}

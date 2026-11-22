package handler

import (
	"net/http"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/port"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type Handler struct {
	PaddleBillingSrv    port.PaddleBillingSrv
	AuthApiMdw          middlewarePort.Authenticator
	IsEnabled           bool
	PaddleWebhookSecret string
}

// NewPaddleBillingHandler registers Billing-owned Paddle transaction,
// subscription-access, and payment webhook routes.
func NewPaddleBillingHandler(mux *http.ServeMux, billingSrv port.PaddleBillingSrv, auth middlewarePort.Authenticator, webhookSecret string, enabled bool) *Handler {
	h := &Handler{
		PaddleBillingSrv:    billingSrv,
		AuthApiMdw:          auth,
		IsEnabled:           enabled,
		PaddleWebhookSecret: webhookSecret,
	}
	mux.HandleFunc("GET /billing/paddle/transactions/{paddleTransactionID}", h.requirePaddleEnabled(h.GetTransaction))
	mux.Handle("GET /billing/paddle/subscription/access", auth.AccessTokenMdw(h.requirePaddleEnabled(h.GetSubscriptionAccess)))
	mux.Handle("POST /billing/paddle/checkout", auth.AccessTokenMdw(h.requirePaddleEnabled(h.CreateCheckout)))
	var webhookHandler http.Handler = http.HandlerFunc(h.PaddleWebhook)
	if enabled && webhookSecret != "" {
		webhookHandler = paddleSDK.NewWebhookVerifier(webhookSecret, paddleSDK.VerifierWithTimestampTolerance(5*time.Second)).Middleware(webhookHandler)
	}
	mux.Handle("POST /billing/paddle/webhook", webhookHandler)
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

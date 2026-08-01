package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req core.CreateProviderReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	provider := &domain.PaymentProvider{
		Name: req.ProviderName,
	}

	switch req.ProviderName {

	case domain.PProviderNameStripe:

		provider.StripeConfig = &domain.StripeConfig{
			PublishableKey: req.StripeConfigReq.PublishableKey,
			SecretKey:      req.StripeConfigReq.SecretKey,
			WebhookSecret:  req.StripeConfigReq.WebhookSecret,
		}

	case domain.PProviderNameLemonSqueezy:

		provider.LemonSqueezyConfig = &domain.LemonSqueezyConfig{
			ApiKey:        req.LemonSqueezyConfigReq.ApiKey,
			StoreID:       req.LemonSqueezyConfigReq.StoreID,
			WebhookSecret: req.LemonSqueezyConfigReq.WebhookSecret,
		}
	}

	if err := h.PProviderSrv.Create(context.Background(), provider); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(provider)
}

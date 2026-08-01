package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func validatePaymentProviderFields(provider *CreateProviderReq) error {
	providerName := domain.PProviderName(
		strings.TrimSpace(string(provider.ProviderName)),
	)

	if providerName == "" {
		return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "provider_name is required")
	}

	switch providerName {
	case domain.PProviderNameStripe:
		if provider.StripeConfigReq == nil {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "stripe config is required")
		}

		if strings.TrimSpace(provider.StripeConfigReq.PublishableKey) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "stripe publishable_key is required")
		}

		if strings.TrimSpace(provider.StripeConfigReq.SecretKey) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "stripe secret_key is required")
		}

		if strings.TrimSpace(provider.StripeConfigReq.WebhookSecret) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "stripe webhook_secret is required")
		}
	case domain.PProviderNameLemonSqueezy:
		if provider.LemonSqueezyConfigReq == nil {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "lemonsqueezy config is required")
		}

		if strings.TrimSpace(provider.LemonSqueezyConfigReq.ApiKey) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "lemonsqueezy api_key is required")
		}

		if strings.TrimSpace(provider.LemonSqueezyConfigReq.StoreID) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "lemonsqueezy store_id is required")
		}

		if strings.TrimSpace(provider.LemonSqueezyConfigReq.WebhookSecret) == "" {
			return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "lemonsqueezy webhook_secret is required")
		}
	default:
		return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "provider_name is invalid")
	}

	return nil
}

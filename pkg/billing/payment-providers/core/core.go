package core

import "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"

type CreateProviderReq struct {
	ProviderName          domain.PProviderName   `json:"provider_name"`
	StripeConfigReq       *StripeConfigReq       `json:"stripe,omitempty"`
	LemonSqueezyConfigReq *LemonSqueezyConfigReq `json:"lemonsqueezy,omitempty"`
}

type LemonSqueezyConfigReq struct {
	ApiKey        string `json:"api_key"`
	StoreID       string `json:"store_id"`
	WebhookSecret string `json:"webhook_secret"`
}

type StripeConfigReq struct {
	PublishableKey string `json:"publishable_key"`
	SecretKey      string `json:"secret_key"`
	WebhookSecret  string `json:"webhook_secret"`
}

func (p *CreateProviderReq) Validate() error {
	return validatePaymentProviderFields(p)
}

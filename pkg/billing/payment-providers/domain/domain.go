package domain

import "time"

type PaymentProvider struct {
	ID                 string              `json:"id"`
	Name               PProviderName       `json:"provider_name"`
	Enabled            bool                `json:"enabled"`
	StripeConfig       *StripeConfig       `json:"stripe,omitempty"`
	LemonSqueezyConfig *LemonSqueezyConfig `json:"lemonsqueezy,omitempty"`
	CreatedAt          *time.Time          `json:"created_at"`
	UpdatedAt          *time.Time          `json:"updated_at"`
}

type StripeConfig struct {
	PublishableKey         string `json:"publishable_key,omitempty"`
	SecretKey              string `json:"-"`
	SecretKeyEncrypted     string `json:"-"`
	WebhookSecret          string `json:"-"`
	WebhookSecretEncrypted string `json:"-"`
}

type LemonSqueezyConfig struct {
	ApiKey                 string `json:"-"`
	ApiKeyEncrypted        string `json:"-"`
	StoreID                string `json:"store_id,omitempty"`
	WebhookSecret          string `json:"-"`
	WebhookSecretEncrypted string `json:"-"`
}

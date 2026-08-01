package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaymentProviderNoSqlModel struct {
	ID      bson.ObjectID        `bson:"_id"`
	Name    domain.PProviderName `bson:"name"`
	Enabled bool                 `bson:"enabled"`

	StripeConfigModel       *StripeConfigModel       `bson:"stripe,omitempty"`
	LemonSqueezyConfigModel *LemonSqueezyConfigModel `bson:"lemonsqueezy,omitempty"`

	CreatedAt *time.Time `bson:"created_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

type StripeConfigModel struct {
	PublishableKey         string `bson:"publishable_key"`
	SecretKeyEncrypted     string `bson:"secret_key_encrypted"`
	WebhookSecretEncrypted string `bson:"webhook_secret_encrypted"`
}

type LemonSqueezyConfigModel struct {
	ApiKeyEncrypted        string `bson:"api_key_encrypted"`
	StoreID                string `bson:"store_id"`
	WebhookSecretEncrypted string `bson:"webhook_secret_encrypted"`
}

func (m *PaymentProviderNoSqlModel) ToDomain(p *domain.PaymentProvider) {
	p.ID = m.ID.Hex()
	p.Name = m.Name
	p.Enabled = m.Enabled
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt

	if m.StripeConfigModel != nil {

		p.StripeConfig = &domain.StripeConfig{
			PublishableKey:         m.StripeConfigModel.PublishableKey,
			SecretKeyEncrypted:     m.StripeConfigModel.SecretKeyEncrypted,
			WebhookSecretEncrypted: m.StripeConfigModel.WebhookSecretEncrypted,
		}
	}

	if m.LemonSqueezyConfigModel != nil {

		p.LemonSqueezyConfig = &domain.LemonSqueezyConfig{
			ApiKeyEncrypted:        m.LemonSqueezyConfigModel.ApiKeyEncrypted,
			StoreID:                m.LemonSqueezyConfigModel.StoreID,
			WebhookSecretEncrypted: m.LemonSqueezyConfigModel.WebhookSecretEncrypted,
		}
	}
}

func FromDomain(p *domain.PaymentProvider,id bson.ObjectID) *PaymentProviderNoSqlModel {
	if p == nil {
		return nil
	}

	model := &PaymentProviderNoSqlModel{
		ID:        id,
		Name:      p.Name,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.StripeConfig != nil {
		model.StripeConfigModel = &StripeConfigModel{
			PublishableKey:         p.StripeConfig.PublishableKey,
			SecretKeyEncrypted:     p.StripeConfig.SecretKeyEncrypted,
			WebhookSecretEncrypted: p.StripeConfig.WebhookSecretEncrypted,
		}
	}

	if p.LemonSqueezyConfig != nil {
		model.LemonSqueezyConfigModel = &LemonSqueezyConfigModel{
			ApiKeyEncrypted:        p.LemonSqueezyConfig.ApiKeyEncrypted,
			StoreID:                p.LemonSqueezyConfig.StoreID,
			WebhookSecretEncrypted: p.LemonSqueezyConfig.WebhookSecretEncrypted,
		}
	}

	return model
}

package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Create(ctx context.Context, pProvider *domain.PaymentProvider) error {
	now := time.Now().UTC()
	pProvider.CreatedAt = &now
	pProvider.UpdatedAt = &now
	pProvider.Enabled = false

	switch pProvider.Name {
	case domain.PProviderNameStripe:
		err := validateStripeCredentials(pProvider.StripeConfig.SecretKey)
		if err != nil {
			return sharedD.ManageError(err,"stripe validation failed")
		}

		secretKey, err := s.EncryptorSrv.Encrypt(pProvider.StripeConfig.SecretKey)
		if err != nil {
			return sharedD.ManageError(err, "error encrypting secret_key")
		}

		webhookSecret, err := s.EncryptorSrv.Encrypt(pProvider.StripeConfig.WebhookSecret)
		if err != nil {
			return sharedD.ManageError(err, "error encrypting webhook_secret")
		}

		pProvider.StripeConfig.SecretKeyEncrypted = secretKey
		pProvider.StripeConfig.WebhookSecretEncrypted = webhookSecret

		pProvider.StripeConfig.SecretKey = ""
		pProvider.StripeConfig.WebhookSecret = ""

	case domain.PProviderNameLemonSqueezy:
		err := validateLemonSqueezyCredentials(pProvider.LemonSqueezyConfig.ApiKey)
		if err != nil {
			return sharedD.ManageError(err,"lemonsqueezy validation failed")
		}

		apiKey, err := s.EncryptorSrv.Encrypt(pProvider.LemonSqueezyConfig.ApiKey)
		if err != nil {
			return sharedD.ManageError(err, "error encrypting api_key")
		}

		webhookSecret, err := s.EncryptorSrv.Encrypt(pProvider.LemonSqueezyConfig.WebhookSecret)
		if err != nil {
			return sharedD.ManageError(err, "error encrypting webhook_secret")
		}

		pProvider.LemonSqueezyConfig.ApiKeyEncrypted = apiKey
		pProvider.LemonSqueezyConfig.WebhookSecretEncrypted = webhookSecret

		pProvider.LemonSqueezyConfig.ApiKey = ""
		pProvider.LemonSqueezyConfig.WebhookSecret = ""
	}

	err := s.PProviderRepo.Insert(ctx, pProvider)
	if err != nil {
		return sharedD.ManageError(err, "")
	}

	return nil
}

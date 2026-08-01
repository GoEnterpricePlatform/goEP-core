package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/domain"
)

type PaymentProviderSrv interface {
	Create(ctx context.Context, pProvider *domain.PaymentProvider) error
}

type PaymentProviderRepo interface {
	Insert(ctx context.Context, pProvider *domain.PaymentProvider) error
}
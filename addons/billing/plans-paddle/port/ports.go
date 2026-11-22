package port

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/domain"
)

type PaddleTransactionRepo interface {
	Insert(ctx context.Context, transaction *domain.PaddleTransaction) error
	FindByPaddleTransactionID(ctx context.Context, paddleTransactionID string) (*domain.PaddleTransaction, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*domain.PaddleTransaction, error)
	FindAllLinkedToSubscriptions(ctx context.Context) ([]*domain.PaddleTransaction, error)
	UpsertFromWebhook(ctx context.Context, transaction *domain.PaddleTransaction, eventAt time.Time) error
}

type PaddleSubscriptionRepo interface {
	Upsert(ctx context.Context, subscription *domain.PaddleSubscription, lastEventAt time.Time) error
	FindByID(ctx context.Context, id string) (*domain.PaddleSubscription, error)
	FindByPaddleSubscriptionID(ctx context.Context, paddleSubscriptionID string) (*domain.PaddleSubscription, error)
	FindCurrentByUserID(ctx context.Context, userID string) (*domain.PaddleSubscription, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*domain.PaddleSubscription, error)
	FindAll(ctx context.Context) ([]*domain.PaddleSubscription, error)
}

type PaddleCustomerRepo interface {
	Upsert(ctx context.Context, customer *domain.PaddleCustomer) error
	FindByUserID(ctx context.Context, userID string) (*domain.PaddleCustomer, error)
	FindByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.PaddleCustomer, error)
}

type PaddleWebhookEventRepo interface {
	StartProcessing(ctx context.Context, event *domain.PaddleWebhookEvent, now time.Time) (alreadyProcessed bool, err error)
	MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error
	MarkFailed(ctx context.Context, eventID string, failedAt time.Time, processingErr error) error
	CreateIndexes() error
}

// PaddleBillingSrv is the service contract Billing exposes to other modules.
type PaddleBillingSrv interface {
	CreateCheckout(ctx context.Context, userID, paddlePlanID, paddlePlanItemID string) (*domain.PaddleTransaction, error)
	GetTransaction(ctx context.Context, paddleTransactionID string) (*domain.PaddleTransaction, error)
	GetSubscriptionAccess(ctx context.Context, userID string) (*domain.PaddleSubscriptionAccess, error)
	HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error
}

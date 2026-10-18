package port

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
)

type PaddlePlanSrv interface {
	Create(ctx context.Context, paddlePlan *domain.PaddlePlan) error
	Get(ctx context.Context, id string) (*domain.PaddlePlan, error)
	GetAll(ctx context.Context, limit, page int64) ([]*domain.PaddlePlan, int64, int64, error)
	Update(ctx context.Context, id string, paddlePlan *domain.PaddlePlan) error
	Delete(ctx context.Context, id string) error
	CreateCheckout(ctx context.Context, userID, planID, itemID string) (*domain.PaddleTransaction, error)
	GetTransaction(ctx context.Context, paddleTransactionID string) (*domain.PaddleTransaction, error)
	HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error
	GetSubscriptionAccess(ctx context.Context, userID string) (*domain.PaddleSubscriptionAccess, error)
}

type PaddlePlanRepo interface {
	Insert(ctx context.Context, paddlePlan *domain.PaddlePlan) error
	Find(ctx context.Context, id string) (*domain.PaddlePlan, error)
	FindAll(ctx context.Context, limit, page int64) ([]*domain.PaddlePlan, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id string, paddlePlan *domain.PaddlePlan) error
	Delete(ctx context.Context, id string) error
}

type PaddlePlanTx interface {
	Create(ctx context.Context, plan *planD.Plan, paddlePlan *domain.PaddlePlan) error
	Update(ctx context.Context, plan *planD.Plan, paddlePlan *domain.PaddlePlan) error
	Delete(ctx context.Context, id string) error
}

type PaddlePlanFileStg interface {
	GetImage(ctx context.Context, imgPath string) (string, error)
}

type PaddleTransactionRepo interface {
	Insert(ctx context.Context, transaction *domain.PaddleTransaction) error
	FindByPaddleTransactionID(ctx context.Context, paddleTransactionID string) (*domain.PaddleTransaction, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*domain.PaddleTransaction, error)
	FindAllLinkedToSubscriptions(ctx context.Context) ([]*domain.PaddleTransaction, error)
	UpsertFromWebhook(ctx context.Context, checkout *domain.PaddleTransaction, eventAt time.Time) error
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

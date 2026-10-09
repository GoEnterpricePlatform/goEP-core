package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

type PaddlePlanSrv interface {
	Create(ctx context.Context, paddlePlan *domain.PaddlePlan) error
	Get(ctx context.Context, id string) (*domain.PaddlePlan, error)
	GetAll(ctx context.Context, limit, page int64) ([]*domain.PaddlePlan, int64, int64, error)
	Update(ctx context.Context, id string, paddlePlan *domain.PaddlePlan) error
	Delete(ctx context.Context, id string) error
	CreateCheckout(ctx context.Context, planID, itemID string) (*domain.PaddleCheckout, error)
	GetCheckout(ctx context.Context, transactionID string) (*domain.PaddleCheckout, error)
	HandlePaddleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error
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

type PaddleCheckoutRepo interface {
	Insert(ctx context.Context, checkout *domain.PaddleCheckout) error
	Find(ctx context.Context, transactionID string) (*domain.PaddleCheckout, error)
	UpdateStatus(ctx context.Context, transactionID, status string, subscriptionID *string, updatedAt string) error
}

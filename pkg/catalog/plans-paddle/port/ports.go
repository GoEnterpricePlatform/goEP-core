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

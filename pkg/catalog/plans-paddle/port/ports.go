package port

import (
	"context"

	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)


type PaddlePlanSrv interface {
	Create(ctx context.Context, paddlePlan *domain.PaddlePlan) error
	Get(ctx context.Context, id string) (*domain.PaddlePlan, error)
}

type PaddlePlanRepo interface {
	Insert(ctx context.Context, paddlePlan *domain.PaddlePlan) error
	Find(ctx context.Context, id string) (*domain.PaddlePlan, error)
}

type PaddlePlanTx interface {
	Create(ctx context.Context, plan *planD.Plan, paddlePlan *domain.PaddlePlan) error
}

type PaddlePlanFileStg interface {
	GetImage(ctx context.Context, imgPath string) (string, error)
}

package port

import (
	"context"
	"io"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

type PlanRepo interface {
	Insert(ctx context.Context, plan *domain.Plan) error
	Find(ctx context.Context, id string) (*domain.Plan, error)
	FindAll(ctx context.Context, limit int64, page int64) ([]*domain.Plan, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id string, plan *domain.Plan) error
	Delete(ctx context.Context, id string) error
}

type PlanSrv interface {
	Get(ctx context.Context, id string) (*domain.Plan, error)
	GetAll(ctx context.Context, limit int64, page int64) ([]*domain.Plan, int64, int64, error)
	Create(ctx context.Context, product *domain.Plan) error
	Update(ctx context.Context, id string, plan *domain.Plan) error
	Patch(ctx context.Context, id string, plan *domain.Plan) (*domain.Plan, error)
	Delete(ctx context.Context, id string) error
}

type PlanFileStg interface {
	GetImage(ctx context.Context, imgPath string) (string, error)
	UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error
}

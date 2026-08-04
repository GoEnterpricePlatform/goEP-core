package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
)

type VariationRepo interface {
	Insert(ctx context.Context, variation *domain.Variation) error
	FindAllWithOptions(ctx context.Context) ([]*domain.Variation, error)
	Update(ctx context.Context, variation *domain.Variation) error
	Delete(ctx context.Context, id string) error
}

type VarOptionRepo interface {
	Insert(ctx context.Context, varOption *domain.VarOption) error
	Update(ctx context.Context, varOption *domain.VarOption) error
	Delete(ctx context.Context, id string, variationID string) error
	FindByIDs(ctx context.Context,ids []string) ([]*domain.VarOption,error)
}

type VariationSrv interface {
	CreateVariation(ctx context.Context, variation *domain.Variation) error
	GetAllVariationsWithOptions(ctx context.Context) ([]*domain.Variation, error)
	UpdateVariation(ctx context.Context, variation *domain.Variation) error
	DeleteVariation(ctx context.Context, id string) error

	CreateVarOption(ctx context.Context, varOption *domain.VarOption) error
	UpdateVarOption(ctx context.Context, varOption *domain.VarOption) error
	DeleteVarOption(ctx context.Context, id string, variationID string) error
}
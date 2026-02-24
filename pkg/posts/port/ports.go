package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
)

type PostRepo interface {
	Insert(ctx context.Context, post *domain.Post) error
	Get(ctx context.Context, id string) (*domain.Post, error)
	FindAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, id string, post *domain.Post) error
	Delete(ctx context.Context, id string) error
}

type PostSrv interface {
	Create(ctx context.Context, post *domain.Post) error
	GetAll(ctx context.Context) ([]*domain.Post, error)
	Update(ctx context.Context, id string, post *domain.Post) error
	Delete(ctx context.Context, id string) error
	Patch(ctx context.Context, id string, post *domain.Post) (*domain.Post, error)
}

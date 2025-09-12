package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/app/users/domain"
)

type UserRepo interface {
	Insert(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

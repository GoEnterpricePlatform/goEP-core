package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/users/domain"
)

type AdminSrv interface {
	ExistsAdmin(ctx context.Context) (bool, error)
	SignUp(ctx context.Context, user *domain.User) error
}
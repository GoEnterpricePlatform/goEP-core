package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/app/roles/domain"
)

type RoleRepo interface {
	Insert(ctx context.Context, role *domain.Role) error
	FindIDs(ctx context.Context, names []string) ([]string, error)
	FindByIDs(ctx context.Context, roleIDs []string) ([]string, error)
	Exists(ctx context.Context, name string) (bool, error)
}

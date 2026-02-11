package port

import (
	"context"

	"github.com/amorindev/go-cms-tmpl/pkg/features/roles/domain"
)

type RoleRepo interface {
	Insert(ctx context.Context, role *domain.Role) error
	Exists(ctx context.Context, name string) (bool, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	AssignPermissions(ctx context.Context, name string, permissionIDs []string) error
}

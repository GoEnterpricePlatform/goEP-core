package port

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
)

type RoleRepo interface {
	Insert(ctx context.Context, role *domain.Role) error
	Exists(ctx context.Context, name string) (bool, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	AssignPermissions(ctx context.Context, name string, permissionIDs []string) error
	FindByIDs(ctx context.Context, roleIDs []string) ([]*domain.Role, error)
	Find(ctx context.Context, id string) (*domain.Role, error)
	FindAll(ctx context.Context) ([]*domain.Role, error)
	Update(ctx context.Context, id string, role *domain.Role) error
	Delete(ctx context.Context, id string) error
}

type RoleSrv interface {
	Create(ctx context.Context, role *domain.Role) error
	Get(ctx context.Context, id string) (*domain.Role, error)
	GetAll(ctx context.Context) ([]*domain.Role, error)
	Update(ctx context.Context, id string, role *domain.Role) error
	Delete(ctx context.Context, id string) error
}

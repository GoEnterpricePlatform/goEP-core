package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/permissions/domain"
)

// PermissionRepo manages system-defined permissions.
//
// domain.PermissionName is used instead of string because permissions are
// NOT dynamic. They must be predefined and controlled by the system,
// unlike roles, which can be created dynamically from the UI.
type PermissionRepo interface {
	Insert(ctx context.Context, permission *domain.Permission) error
	Exists(ctx context.Context, name domain.PermissionName) (bool, error)
	FindByNames(ctx context.Context, names []domain.PermissionName) ([]*domain.Permission, error)
}

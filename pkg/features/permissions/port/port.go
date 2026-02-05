package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/permissions/domain"
)

type PermissionRepo interface {
	Insert(ctx context.Context, permission *domain.Permission) error
	Exists(ctx context.Context, name string) (bool, error)
}

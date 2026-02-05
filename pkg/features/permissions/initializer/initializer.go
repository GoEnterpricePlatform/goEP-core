package initializer

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/permissions/domain"
	"github.com/amorindev/go-tmpl/pkg/features/permissions/port"
)

type Initializer struct {
	PermissionRepo port.PermissionRepo
}

func NewPermissionItz(permissionRepo port.PermissionRepo) *Initializer {
	return &Initializer{
		PermissionRepo: permissionRepo,
	}
}

func (i *Initializer) SeedEssentialPermissions(ctx context.Context) ([]*domain.Permission, error) {
	permissions := []*domain.Permission{
		domain.NewPermission(domain.PAdminAccess),
		domain.NewPermission(domain.PAdminManage),
		domain.NewPermission(domain.PUserRead),
		domain.NewPermission(domain.PUserChangeActive),
		domain.NewPermission(domain.PSettingsRead),
		domain.NewPermission(domain.PSettingsUpdate),
	}

	for _, permission := range permissions {
		exists, err := i.PermissionRepo.Exists(ctx, permission.Name)
		if err != nil {
			return nil, err
		}
		if !exists {
			if err := i.PermissionRepo.Insert(ctx, permission); err != nil {
				return nil, err
			}
		}
	}
	return permissions, nil
}

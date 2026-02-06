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
	permissionNames := []domain.PermissionName{
		domain.PAdminAccess,
		domain.PAdminManage,
		domain.PUserRead,
		domain.PUserChangeActive,
		domain.PSettingsRead,
		domain.PSettingsUpdate,
	}

	for _, name := range permissionNames {
		exists, err := i.PermissionRepo.Exists(ctx, name)
		if err != nil {
			return nil, err
		}
		if !exists {
			if err := i.PermissionRepo.Insert(ctx, domain.NewPermission(name)); err != nil {
				return nil, err
			}
		}
	}

	permissions, err := i.PermissionRepo.FindByNames(ctx, permissionNames)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

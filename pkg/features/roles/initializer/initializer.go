package initializer

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/roles/domain"
	"github.com/amorindev/go-tmpl/pkg/features/roles/port"
)

type Initializer struct {
	RoleRepo port.RoleRepo
}

func NewRoleItz(roleRepo port.RoleRepo) *Initializer {
	return &Initializer{
		RoleRepo: roleRepo,
	}
}

func (i *Initializer) SeedEssentialRoles(ctx context.Context) error {
	roles := []*domain.Role{
		domain.NewRole(string(domain.RoleAdmin)),
		domain.NewRole(string(domain.RoleUser)),
	}

	for _, role := range roles {
		exists, err := i.RoleRepo.Exists(ctx, role.Name)
		if err != nil {
			return err
		}
		if !exists {
			if err := i.RoleRepo.Insert(ctx, role); err != nil {
				return err
			}
		}
	}
	return nil
}

package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) Update(ctx context.Context, id string, role *domain.Role) error {
	err := s.RoleRepo.Update(ctx, id, role)
	if err != nil {
		return sharedD.ManageError(err, "error updating role")
	}
	return nil
}

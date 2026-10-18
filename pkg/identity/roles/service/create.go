package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) Create(ctx context.Context, role *domain.Role) error {
	err := s.RoleRepo.Insert(ctx, role)
	if err != nil {
		return sharedD.ManageError(err, "error creating role")
	}

	return nil
}

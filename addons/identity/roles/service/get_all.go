package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) GetAll(ctx context.Context) ([]*domain.Role, error) {
	roles, err := s.RoleRepo.FindAll(ctx)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting all roles")
	}
	return roles, nil
}


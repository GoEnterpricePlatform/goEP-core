package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) Get(ctx context.Context, id string) (*domain.Role, error) {
	role, err := s.RoleRepo.Find(ctx, id)
	if err != nil {
		return nil, sharedD.ManageError(err, "error getting role")
	}
	return role, nil
}


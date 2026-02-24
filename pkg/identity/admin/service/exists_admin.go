package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) ExistsAdmin(ctx context.Context) (bool, error) {
	exists, err := s.UserRepo.ExistsAdmin(ctx)
	if err != nil {
		return false, domain.ManageError(err, "checking admin existence")
	}

	return exists, nil
}

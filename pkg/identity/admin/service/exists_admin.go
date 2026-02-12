package service

import (
	"context"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) ExistsAdmin(ctx context.Context) (bool, error) {
	exists, err := s.UserRepo.ExistsAdmin(ctx)
	if err != nil {
		return false, domain.ManageError(err, "checking admin existence")
	}

	return exists, nil
}

package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) CreateVarOption(ctx context.Context, varOption *domain.VarOption) error {
	now := time.Now().UTC()
	varOption.CreatedAt = &now
	varOption.UpdatedAt = &now

	err := s.VarOptionRepo.Insert(ctx,varOption)
    if err != nil {
      return sharedD.ManageError(err, "error creating varOption")
    }

	return nil
}

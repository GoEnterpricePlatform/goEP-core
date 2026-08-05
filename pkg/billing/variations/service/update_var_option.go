package service

import (
	"context"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) UpdateVarOption(ctx context.Context,  varOption *domain.VarOption) error {
	now := time.Now().UTC()
	varOption.UpdatedAt = &now

	err := s.VarOptionRepo.Update(ctx, varOption)
	if err != nil {
		return sharedD.ManageError(err, "error updating varOption")
	}
    return nil
}

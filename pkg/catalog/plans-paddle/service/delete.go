package service

import (
	"context"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.PaddlePlanTx.Delete(ctx, id); err != nil {
		return sharedD.ManageError(err, "error deleting paddle plan")
	}
	return nil
}

package service

import (
	"context"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)


func (s *Service) GetCheckout(ctx context.Context, transactionID string) (*paddlePlanD.PaddleCheckout, error) {
	checkout, err := s.PaddleCheckoutRepo.Find(ctx, transactionID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding Paddle checkout")
	}
	return checkout, nil
}
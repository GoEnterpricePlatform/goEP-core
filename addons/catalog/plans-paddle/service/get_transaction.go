package service

import (
	"context"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) GetTransaction(ctx context.Context, paddleTransactionID string) (*paddlePlanD.PaddleTransaction, error) {
	transaction, err := s.PaddleTransactionRepo.FindByPaddleTransactionID(ctx, paddleTransactionID)
	if err != nil {
		return nil, sharedD.ManageError(err, "error finding Paddle transaction")
	}
	return transaction, nil
}

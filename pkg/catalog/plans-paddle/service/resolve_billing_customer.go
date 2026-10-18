package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) resolveBillingCustomer(ctx context.Context, userID, paddleCustomerID string) (string, string, error) {
	if s.PaddleCustomerRepo == nil {
		return userID, "", nil
	}
	if paddleCustomerID != "" && userID == "" {
		customer, err := s.PaddleCustomerRepo.FindByPaddleCustomerID(ctx, paddleCustomerID)
		if err == nil {
			userID = customer.UserID
		} else if !errors.Is(err, sharedD.ErrNotFound) {
			return userID, "", fmt.Errorf("error finding GoEP user for Paddle customer: %w", err)
		}
	}
	if userID == "" {
		return userID, "", nil
	}
	if paddleCustomerID != "" {
		customer := &domain.PaddleCustomer{UserID: userID, PaddleCustomerID: paddleCustomerID}
		if err := s.PaddleCustomerRepo.Upsert(ctx, customer); err != nil {
			return userID, "", fmt.Errorf("error saving Paddle billing identity: %w", err)
		}
		return userID, customer.ID, nil
	}
	if customer, err := s.PaddleCustomerRepo.FindByUserID(ctx, userID); err == nil {
		return userID, customer.ID, nil
	} else if !errors.Is(err, sharedD.ErrNotFound) {
		return userID, "", fmt.Errorf("error finding Paddle billing identity: %w", err)
	}
	return userID, "", nil
}

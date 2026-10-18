package service

import (
	"context"
	"errors"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

func (s *Service) GetSubscriptionAccess(ctx context.Context, userID string) (*domain.PaddleSubscriptionAccess, error) {
	if userID == "" {
		return nil, sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required")
	}
	subscription, err := s.PaddleSubscriptionRepo.FindCurrentByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sharedD.ErrNotFound) {
			return &domain.PaddleSubscriptionAccess{HasAccess: false}, nil
		}
		return nil, sharedD.ManageError(err, "error finding current Paddle subscription")
	}
	return &domain.PaddleSubscriptionAccess{HasAccess: true, Subscription: subscription}, nil
}

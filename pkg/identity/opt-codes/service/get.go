package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/domain"
)

func (s *Service) Get(ctx context.Context, id, userID string) (*domain.OtpCode, error) {
	return s.OtpCodeRepo.Find(ctx, id, userID)
}

package service

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/opt-codes/domain"
)

func (s *Service) Get(ctx context.Context, id, userID string) (*domain.OtpCode, error) {
	return s.OtpCodeRepo.Find(ctx, id, userID)
}

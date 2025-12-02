package service

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/session/domain"
)

func (s *Service) GetByRTokenID(ctx context.Context, rTokenID, userID string) (*domain.Session, error) {
	session, err := s.SessionRepo.FindByRTokenID(ctx, rTokenID, userID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

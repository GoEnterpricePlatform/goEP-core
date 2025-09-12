package service

import (
	"context"
)

func (s *Service) DeleteByRTokenID(ctx context.Context, rTokenID string) error {
	err := s.SessionRepo.DeleteByRTokenID(ctx, rTokenID)
	if err != nil {
		return err
	}
	return nil
}

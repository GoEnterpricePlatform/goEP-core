package service

import (
	"context"

	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) SignOut(ctx context.Context, rTokenID string) error {
	err := s.SessionSrv.DeleteByRTokenID(ctx, rTokenID)
	if err != nil {
		return sharedD.ManageError(err, "failed to sign out")
	}
	return nil
}

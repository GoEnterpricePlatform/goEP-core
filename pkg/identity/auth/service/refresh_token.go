package service

import (
	"context"

	sessionD "github.com/amorindev/go-cms-tmpl/pkg/identity/session/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)


func (s *Service) RefreshToken(ctx context.Context, rTokenID string, userID string) (*sessionD.Session, error) {
	session, err := s.SessionSrv.GetByRTokenID(ctx, rTokenID, userID)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	if session.Revoked {
		return nil, sharedD.ManageError(err, "")
	}

	err = s.SessionSrv.DeleteByRTokenID(ctx, rTokenID)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	user, err := s.UserRepo.Find(ctx, userID)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	// Create session
	newSession := sessionD.NewSession(user.ID, session.RememberMe)

	err = s.SessionSrv.Create(ctx, newSession, nil, user.Email)
	if err != nil {
		return nil, sharedD.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return newSession, nil
}

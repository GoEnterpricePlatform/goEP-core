package service

import (
	"context"

	"github.com/amorindev/go-tmpl/internal/encryption"
	sessionD "github.com/amorindev/go-tmpl/pkg/app/session/domain"
	userD "github.com/amorindev/go-tmpl/pkg/app/users/domain"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

func (s *Service) SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, error) {
	// Verify if user exists
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, dShared.ManageError(err, "")
	}

	if !user.IsActive {
		return nil, nil, dShared.ManageError(dShared.ErrAccountInactive, "")
	}

	err = encryption.CheckPassword(password, user.UserPassAuth.PasswordHash)
	if err != nil {
		return nil, nil, dShared.ManageError(err, "")
	}

	// Create session
	session := sessionD.NewSession(user.ID.(string), rememberMe)

	err = s.SessionSrv.Create(ctx, session, nil, email)
	if err != nil {
		return nil, nil, dShared.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return user, session, nil
}

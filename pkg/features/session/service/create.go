package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
)

func (s *Service) Create(ctx context.Context, session *domain.Session, roles []string, email string) error {
	// Create access token
	aToken, aTokenExpIn, err := s.TokenSrv.CreateAccessToken(session.UserID.(string), email, roles)
	if err != nil {
		return err
	}

	// Create refresh token
	rTokenID, rToken, rTokenExpIn, err := s.TokenSrv.CreateRefreshToken(session.UserID.(string), session.RememberMe)
	if err != nil {
		return err
	}

	// Create session
	now := time.Now().UTC()
	session.AccessToken = aToken
	session.AccessTokenExpIn = aTokenExpIn
	session.RefreshTokenID = rTokenID
	session.RefreshToken = rToken
	session.RefreshTokenExpIn = rTokenExpIn
	session.Revoked = false
	session.CreatedAt = &now

	err = s.SessionRepo.Insert(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

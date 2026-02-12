package service

import (
	"context"

	sessionD "github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
	userD "github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) GetSession(ctx context.Context, rTokenID string, userID string) (*userD.User,*sessionD.Session, error) {
	session, err := s.SessionSrv.GetByRTokenID(ctx, rTokenID, userID)
	if err != nil {
		return nil,nil, sharedD.ManageError(err, "")
	}

	if session.Revoked {
		return nil,nil, sharedD.ManageError(err, "")
	}

	user, err := s.UserRepo.Find(ctx, userID)
	if err != nil {
		return nil,nil, sharedD.ManageError(err, "")
	}
    // If the user has an image path, retrieve the image URL and set it
	if user.ImgPath != nil {
		url, err := s.UserFileStg.GetImage(ctx, *user.ImgPath)
		if err != nil {
			return nil, nil, sharedD.ManageError(err, "")
		}
		user.ImgUrl = &url
	}

	// Create session
	newSession := sessionD.NewSession(user.ID, session.RememberMe)

	err = s.SessionSrv.Create(ctx, newSession, nil, user.Email)
	if err != nil {
		return nil, nil,sharedD.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return user,newSession, nil
}

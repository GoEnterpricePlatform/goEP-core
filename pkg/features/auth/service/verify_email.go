package service

import (
	"context"
	"time"

	otpCodeD "github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	sessionD "github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) VerifyEmail(ctx context.Context, otpID string, otpCode string, userID string) (*domain.User, *sessionD.Session, error) {
	otp, err := s.OtpCodeSrv.Get(ctx, otpID, userID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed getting otpCode")
	}

	if time.Now().After(*otp.ExpiresAt) {
		return nil, nil, sharedD.ManageError(sharedD.ErrOtpExpired, "OTP code has expired")
	}

	if otpCode != otp.OptCode {
		return nil, nil, sharedD.ManageError(sharedD.ErrInvalidOtpCode, "invalid OTP code")
	}

	if otp.Purpose != otpCodeD.VerifyEmailPurpose {
		return nil, nil, sharedD.ManageError(sharedD.ErrOtpPurposeNotAllowed, "OTP not valid for email verification")
	}

	err = s.UserRepo.ConfirmEmail(context.Background(), userID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to confirm user email")
	}

	user, err := s.UserRepo.Find(ctx, userID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to retrieve user data")
	}

	// If the user has an image path, retrieve the image URL and set it
	if user.ImgPath != nil {
		url, err := s.UserFileStg.GetImage(ctx, *user.ImgPath)
		if err != nil {
			return nil, nil, sharedD.ManageError(err, "")
		}
		user.ImgUrl = &url
	}

	session := &sessionD.Session{
		UserID:     user.ID,
		RememberMe: false,
	}

	err = s.SessionSrv.Create(ctx, session, nil, user.Email)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to create user session")
	}

	err = s.OtpCodeSrv.Delete(ctx, otp.ID.(string))
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to delete OTP code")
	}

	user.UserPassAuth = nil

	return user, session, nil
}

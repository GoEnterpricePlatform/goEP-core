package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/internal/encryption"
	otpCodeD "github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	sessionD "github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
	userD "github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	dShared "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	sharedDomain "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, string, error) {
	// Verify if user exists
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, "", dShared.ManageError(err, "")
	}

	if !user.IsActive {
		return nil, nil, "", dShared.ManageError(dShared.ErrAccountInactive, "")
	}

	err = encryption.CheckPassword(password, user.UserPassAuth.PasswordHash)
	if err != nil {
		return nil, nil, "", dShared.ManageError(err, "")
	}

	// If the user has an image path, retrieve the image URL and set it
	if user.ImgPath != nil {
		url, err := s.UserFileStg.GetImage(ctx, *user.ImgPath)
		if err != nil {
			return nil, nil, "", dShared.ManageError(err, "")
		}
		user.ImgUrl = &url
	}

	if !user.EmailVerified {
		now := time.Now().UTC()

		expiresAt := now.Add(time.Hour)

		// Create otp
		otp := &otpCodeD.OtpCode{
			UserID:    user.ID,
			Purpose:   otpCodeD.VerifyEmailPurpose,
			Used:      false,
			ExpiresAt: &expiresAt,
			CreatedAt: &now,
		}

		// Create a new OTP code and save it to the database
		err = s.OtpCodeSrv.Create(ctx, otp)
		if err != nil {
			return nil, nil, "", sharedDomain.ManageError(err, "")
		}

		// Send the verification email to the user with the OTP code
		err = s.MailerSrv.SendVerifyEmail(user.Email, otp.OptCode)
		if err != nil {
			return nil, nil, "", dShared.ManageError(err, "")
		}
		return user, nil, otp.ID.(string), nil
	}

	// Create session
	session := sessionD.NewSession(user.ID.(string), rememberMe)

	err = s.SessionSrv.Create(ctx, session, nil, email)
	if err != nil {
		return nil, nil, "", dShared.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return user, session, "", nil
}

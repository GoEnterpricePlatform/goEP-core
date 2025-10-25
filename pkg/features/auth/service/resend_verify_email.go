package service

import (
	"context"
	"time"

	"github.com/amorindev/go-tmpl/pkg/features/opt-codes/domain"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
	sharedD "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

func (s *Service) ResendVerifyEmail(ctx context.Context, email string) (string, error) {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", dShared.ManageError(err, "")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Hour)

	// Create otp
	otp := &domain.OtpCode{
		UserID:    user.ID,
		Purpose:   domain.VerifyEmailPurpose,
		Used:      false,
		ExpiresAt: &expiresAt,
		CreatedAt: &now,
	}

	// Create a new OTP code and save it to the database
	err = s.OtpCodeSrv.Create(ctx, otp)
	if err != nil {
		return "", sharedD.ManageError(err, "error creating otp code")
	}

	// Send the verification email to the user with the OTP code
	err = s.MailerSrv.SendVerifyEmail(user.Email, otp.OptCode)
	if err != nil {
		return "", sharedD.ManageError(err, "error sending email")
	}

	return otp.ID.(string), nil
}

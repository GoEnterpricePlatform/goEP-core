package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/internal/encryption"
	otpCodeD "github.com/amorindev/go-cms-tmpl/pkg/features/opt-codes/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	sharedDomain "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

// SignUp registers a new user, hashes the password, and saves it to the repository and return otpID
func (s *Service) SignUp(ctx context.Context, user *domain.User) (string, error) {
	// Check if email already exists
	exists, err := s.UserRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return "", sharedDomain.ManageError(err, "checking email existence")
	}

	if exists {
		return "", sharedDomain.ManageError(sharedDomain.ErrDuplicateKey, "email already in use")
	}

	// Hash the user's password
	hashPass, err := encryption.HashPassword(user.UserPassAuth.Password)
	if err != nil {
		return "", sharedDomain.ManageError(err, "error hashing password")
	}
	user.UserPassAuth.PasswordHash = hashPass

	// Create the user
	now := time.Now().UTC()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.IsActive = true
	user.EmailVerified = false
	user.UserPassAuth.CreatedAt = &now
	user.UserPassAuth.UpdatedAt = &now

	// Save user
	err = s.UserRepo.Insert(ctx, user)
	if err != nil {
		return "", sharedDomain.ManageError(err, "error inserting user")
	}

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
		return "", sharedDomain.ManageError(err, "error creating otp code")
	}

	// Send the verification email to the user with the OTP code
	err = s.MailerSrv.SendVerifyEmail(user.Email, otp.OptCode)
	if err != nil {
		return "", sharedDomain.ManageError(err, "error sending email")
	}

	// Clear password hash from memory for security
	user.UserPassAuth = nil

	return otp.ID.(string), nil
}

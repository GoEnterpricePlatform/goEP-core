package service

import (
	"context"
	"time"

	"github.com/amorindev/go-tmpl/internal/encryption"
	"github.com/amorindev/go-tmpl/pkg/features/users/domain"
	sharedDomain "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

// SignUp registers a new user, hashes the password, and saves it to the repository.
func (s *Service) SignUp(ctx context.Context, user *domain.User) error {
	// Check if email already exists
	exists, err := s.UserRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return sharedDomain.ManageError(err, "checking email existence")
	}

	if exists {
		return sharedDomain.ManageError(sharedDomain.ErrDuplicateKey, "email already in use")
	}

	// Hash the user's password
	hashPass, err := encryption.HashPassword(user.UserPassAuth.Password)
	if err != nil {
		return sharedDomain.ManageError(err, "error hashing password")
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
		return sharedDomain.ManageError(err, "error inserting user")
	}

	// Clear password hash from memory for security
	user.UserPassAuth = nil

	return nil
}

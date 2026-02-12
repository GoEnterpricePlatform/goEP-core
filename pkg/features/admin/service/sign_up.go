package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/internal/encryption"
	roleD "github.com/amorindev/go-cms-tmpl/pkg/features/roles/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) SignUp(ctx context.Context, user *domain.User) error {
	// verify if admin exists
	exists, err := s.UserRepo.ExistsAdmin(ctx)
	if err != nil {
		return sharedD.ManageError(err, "checking admin existence")
	}

	// podira ser duplicate key
	if exists {
		return sharedD.ManageError(sharedD.ErrDuplicateKey, "admin already exists")
	}

	// Hash the user's password
	hashPass, err := encryption.HashPassword(user.UserPassAuth.Password)
	if err != nil {
		return sharedD.ManageError(err, "error hashing password")
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

	// assign role admin
	role, err := s.RoleRepo.FindByName(ctx, string(roleD.RoleSystemAdmin))
	if err != nil {
		return sharedD.ManageError(err, "error finding the admin role")
	}
	user.RoleIDs = []string{role.ID.(string)}

	// save admin
	err = s.UserRepo.Insert(ctx, user)
	if err != nil {
		return sharedD.ManageError(err, "error inserting user admin")
	}

	user.UserPassAuth = nil

	return nil
}

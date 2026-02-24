package service

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/encryption"
	sessionD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/domain"
	userD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, error) {
	// Verify if user exists
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, sharedD.ManageError(sharedD.ErrInvalidCredentials, "")
	}

	if !user.IsActive {
		return nil, nil, sharedD.ManageError(sharedD.ErrAccountInactive, "")
	}

	err = encryption.CheckPassword(password, user.UserPassAuth.PasswordHash)
	if err != nil {
		return nil, nil, sharedD.ManageError(sharedD.ErrInvalidCredentials, "")
	}

	// get Roles
	roles, err := s.RoleRepo.FindByIDs(ctx, user.RoleIDs)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "")
	}

	roleNames := make([]string, 0, len(roles))

	// Collect all permission IDs, no duplicates
	permissionIDMap := make(map[string]struct{})

	for _, role := range roles {
		for _, pid := range role.PermissionIDs {
			permissionIDMap[pid] = struct{}{}
		}
		roleNames = append(roleNames, role.Name)

	}

	// From the previous one we get "{"read_users": {},"write_users": {},"export_users": {}}",
	// now we convert it to a slice
	var permissionIDs []string
	for id := range permissionIDMap {
		permissionIDs = append(permissionIDs, id)
	}

	// Get permissions
	permissions, err := s.PermissionRepo.FindByIDs(ctx, permissionIDs)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "")
	}

	permissionNames := make([]string, 0, len(permissions))

	for _, p := range permissions {
		permissionNames = append(permissionNames, p.Name)
	}

	// Create session
	session := sessionD.NewSession(user.ID, rememberMe)

	err = s.SessionSrv.Create(ctx, session, roleNames, permissionNames, email)
	if err != nil {
		return nil, nil, sharedD.ManageError(sharedD.ErrInvalidCredentials, "")
	}

	user.UserPassAuth = nil

	return user, session, nil
}

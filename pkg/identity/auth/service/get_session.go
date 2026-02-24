package service

import (
	"context"

	sessionD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/domain"
	userD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (s *Service) GetSession(ctx context.Context, rTokenID string, userID string) (*userD.User, *sessionD.Session, error) {
	session, err := s.SessionSrv.GetByRTokenID(ctx, rTokenID, userID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "")
	}

	if session.Revoked {
		return nil, nil, sharedD.ManageError(err, "")
	}

	user, err := s.UserRepo.Find(ctx, userID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "")
	}

	// If the user has an image path, retrieve the image URL and set it
	if user.ImgPath != nil {
		url, err := s.UserFileStg.GetImage(ctx, *user.ImgPath)
		if err != nil {
			return nil, nil, sharedD.ManageError(err, "")
		}
		user.ImgUrl = &url
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
	newSession := sessionD.NewSession(user.ID, session.RememberMe)

	err = s.SessionSrv.Create(ctx, newSession, roleNames, permissionNames, user.Email)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return user, newSession, nil
}

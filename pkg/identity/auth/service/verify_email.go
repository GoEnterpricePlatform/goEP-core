package service

import (
	"context"
	"time"

	otpCodeD "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/domain"
	sessionD "github.com/amorindev/go-cms-tmpl/pkg/identity/session/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/domain"
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

	session := &sessionD.Session{
		UserID:     user.ID,
		RememberMe: false,
	}

	err = s.SessionSrv.Create(ctx, session, roleNames, permissionNames, user.Email)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to create user session")
	}

	err = s.OtpCodeSrv.Delete(ctx, otp.ID)
	if err != nil {
		return nil, nil, sharedD.ManageError(err, "failed to delete OTP code")
	}

	user.UserPassAuth = nil

	return user, session, nil
}

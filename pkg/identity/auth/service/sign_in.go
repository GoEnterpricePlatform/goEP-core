package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/encryption"
	otpCodeD "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/domain"
	sessionD "github.com/amorindev/go-cms-tmpl/pkg/identity/session/domain"
	userD "github.com/amorindev/go-cms-tmpl/pkg/identity/users/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	sharedDomain "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (s *Service) SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, string, error) {
	// Verify if user exists
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, "", sharedD.ManageError(err, "")
	}

	if !user.IsActive {
		return nil, nil, "", sharedD.ManageError(sharedD.ErrAccountInactive, "")
	}

	err = encryption.CheckPassword(password, user.UserPassAuth.PasswordHash)
	if err != nil {
		return nil, nil, "", sharedD.ManageError(err, "")
	}

	// If the user has an image path, retrieve the image URL and set it
	if user.ImgPath != nil {
		url, err := s.UserFileStg.GetImage(ctx, *user.ImgPath)
		if err != nil {
			return nil, nil, "", sharedD.ManageError(err, "")
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
			return nil, nil, "", sharedD.ManageError(err, "")
		}
		return user, nil, otp.ID, nil
	}

	// get Roles
	roles, err := s.RoleRepo.FindByIDs(ctx, user.RoleIDs)
	if err != nil {
		return nil, nil, "",sharedD.ManageError(err, "")
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
		return nil, nil, "",sharedD.ManageError(err, "")
	}

	permissionNames := make([]string, 0, len(permissions))

	for _, p := range permissions {
		permissionNames = append(permissionNames, p.Name)
	}

	// Create session
	session := sessionD.NewSession(user.ID, rememberMe)

	err = s.SessionSrv.Create(ctx, session, roleNames,permissionNames, email)
	if err != nil {
		return nil, nil, "", sharedD.ManageError(err, "")
	}

	user.UserPassAuth = nil

	return user, session, "", nil
}

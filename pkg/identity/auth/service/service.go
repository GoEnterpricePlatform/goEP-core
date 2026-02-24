package service

import (
	authP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/port"
	mailerP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/port"
	otpCodeP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/port"
	permissionP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/port"
	sessionP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/port"
	userP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/port"
	roleP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/port"
)

var _ authP.AuthSrv = &Service{}

type Service struct {
	UserRepo       userP.UserRepo
	RoleRepo       roleP.RoleRepo
	PermissionRepo permissionP.PermissionRepo
	UserFileStg    userP.UserFileStg
	SessionSrv     sessionP.SessionSrv
	OtpCodeSrv     otpCodeP.OtpCodeSrv
	MailerSrv      mailerP.MailerSrv
}

func NewAuthSrv(userRepo userP.UserRepo, roleRepo roleP.RoleRepo, permissionRepo permissionP.PermissionRepo, userFileStg userP.UserFileStg, sessionSrv sessionP.SessionSrv, otpCodeSrv otpCodeP.OtpCodeSrv, mailerSrv mailerP.MailerSrv) *Service {
	return &Service{
		UserRepo:       userRepo,
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
		UserFileStg:    userFileStg,
		SessionSrv:     sessionSrv,
		OtpCodeSrv:     otpCodeSrv,
		MailerSrv:      mailerSrv,
	}
}

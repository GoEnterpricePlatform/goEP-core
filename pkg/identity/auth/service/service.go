package service

import (
	authP "github.com/amorindev/go-cms-tmpl/pkg/identity/auth/port"
	mailerP "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/port"
	otpCodeP "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/port"
	permissionP "github.com/amorindev/go-cms-tmpl/pkg/identity/permissions/port"
	roleP "github.com/amorindev/go-cms-tmpl/pkg/identity/roles/port"
	sessionP "github.com/amorindev/go-cms-tmpl/pkg/identity/session/port"
	userP "github.com/amorindev/go-cms-tmpl/pkg/identity/users/port"
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

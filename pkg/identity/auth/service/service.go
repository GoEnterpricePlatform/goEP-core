package service

import (
	authP "github.com/amorindev/go-cms-tmpl/pkg/identity/auth/port"
	mailerP "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/port"
	otpCodeP "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/port"
	sessionP "github.com/amorindev/go-cms-tmpl/pkg/identity/session/port"
	userP "github.com/amorindev/go-cms-tmpl/pkg/identity/users/port"
)

var _ authP.AuthSrv = &Service{}

type Service struct {
	UserRepo    userP.UserRepo
	UserFileStg userP.UserFileStg
	SessionSrv  sessionP.SessionSrv
	OtpCodeSrv  otpCodeP.OtpCodeSrv
	MailerSrv   mailerP.MailerSrv
}

func NewAuthSrv(userRepo userP.UserRepo, userFileStg userP.UserFileStg, sessionSrv sessionP.SessionSrv, otpCodeSrv otpCodeP.OtpCodeSrv, mailerSrv mailerP.MailerSrv) *Service {
	return &Service{
		UserRepo:    userRepo,
		UserFileStg: userFileStg,
		SessionSrv:  sessionSrv,
		OtpCodeSrv:  otpCodeSrv,
		MailerSrv:   mailerSrv,
	}
}

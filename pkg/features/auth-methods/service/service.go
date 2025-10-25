package service

import (
	authMethodP "github.com/amorindev/go-tmpl/pkg/features/auth-methods/port"
	mailerP "github.com/amorindev/go-tmpl/pkg/features/mailer/port"
	otpCodeP "github.com/amorindev/go-tmpl/pkg/features/opt-codes/port"
	sessionP "github.com/amorindev/go-tmpl/pkg/features/session/port"
	userP "github.com/amorindev/go-tmpl/pkg/features/users/port"
)

var _ authMethodP.AuthMethodSrv = &Service{}

type Service struct {
	UserRepo    userP.UserRepo
	UserFileStg userP.UserFileStg
	SessionSrv  sessionP.SessionSrv
	OtpCodeSrv  otpCodeP.OtpCodeSrv
	MailerSrv   mailerP.MailerSrv
}

func NewAuthMethodSrv(userRepo userP.UserRepo, userFileStg userP.UserFileStg, sessionSrv sessionP.SessionSrv, otpCodeSrv otpCodeP.OtpCodeSrv, mailerSrv mailerP.MailerSrv) *Service {
	return &Service{
		UserRepo:    userRepo,
		UserFileStg: userFileStg,
		SessionSrv:  sessionSrv,
		OtpCodeSrv:  otpCodeSrv,
		MailerSrv:   mailerSrv,
	}
}

package service

import (
	authMethodP "github.com/amorindev/go-tmpl/pkg/features/auth-methods/port"
	sessionP "github.com/amorindev/go-tmpl/pkg/features/session/port"
	userP "github.com/amorindev/go-tmpl/pkg/features/users/port"
)

var _ authMethodP.AuthMethodSrv = &Service{}

type Service struct {
	UserRepo    userP.UserRepo
	UserFileStg userP.UserFileStg
	SessionSrv  sessionP.SessionSrv
}

func NewAuthMethodSrv(userRepo userP.UserRepo, userFileStg userP.UserFileStg, sessionSrv sessionP.SessionSrv) *Service {
	return &Service{
		UserRepo:    userRepo,
		UserFileStg: userFileStg,
		SessionSrv:  sessionSrv,
	}
}

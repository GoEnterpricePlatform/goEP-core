package service

import (
	authMethodP "github.com/amorindev/go-tmpl/pkg/app/auth-methods/port"
	sessionP "github.com/amorindev/go-tmpl/pkg/app/session/port"
	userP "github.com/amorindev/go-tmpl/pkg/app/users/port"
)

var _ authMethodP.AuthMethodSrv = &Service{}

type Service struct {
	UserRepo   userP.UserRepo
	SessionSrv sessionP.SessionSrv
}

func NewAuthMethodSrv(userRepo userP.UserRepo, sessionSrv sessionP.SessionSrv) *Service {
	return &Service{
		UserRepo:   userRepo,
		SessionSrv: sessionSrv,
	}
}

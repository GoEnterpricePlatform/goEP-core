package service

import (
	adminP "github.com/amorindev/go-cms-tmpl/pkg/features/admin/port"
	roleP "github.com/amorindev/go-cms-tmpl/pkg/features/roles/port"
	sessionP "github.com/amorindev/go-cms-tmpl/pkg/features/session/port"
	userP "github.com/amorindev/go-cms-tmpl/pkg/features/users/port"
)

var _ adminP.AdminSrv = &Service{}

type Service struct {
	UserRepo   userP.UserRepo
	RoleRepo   roleP.RoleRepo
	SessionSrv sessionP.SessionSrv
}

func NewAdminSrv(userRepo userP.UserRepo, roleRepo roleP.RoleRepo, sessionSrv sessionP.SessionSrv) *Service {
	return &Service{
		UserRepo:   userRepo,
		RoleRepo:   roleRepo,
		SessionSrv: sessionSrv,
	}
}

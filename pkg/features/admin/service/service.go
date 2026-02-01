package service

import (
	adminP "github.com/amorindev/go-tmpl/pkg/features/admin/port"
	roleP "github.com/amorindev/go-tmpl/pkg/features/roles/port"
	userP "github.com/amorindev/go-tmpl/pkg/features/users/port"
)

var _ adminP.AdminSrv = &Service{}

type Service struct {
	UserRepo userP.UserRepo
	RoleRepo roleP.RoleRepo
}

func NewAdminSrv(userRepo userP.UserRepo, roleRepo roleP.RoleRepo) *Service {
	return &Service{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

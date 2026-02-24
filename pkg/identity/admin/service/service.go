package service

import (
	adminP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/port"
	permissionP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/port"
	roleP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/port"
	sessionP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/port"
	userP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/port"
)

var _ adminP.AdminSrv = &Service{}

type Service struct {
	UserRepo       userP.UserRepo
	RoleRepo       roleP.RoleRepo
	PermissionRepo permissionP.PermissionRepo
	SessionSrv     sessionP.SessionSrv
}

func NewAdminSrv(userRepo userP.UserRepo, roleRepo roleP.RoleRepo, permissionRepo permissionP.PermissionRepo, sessionSrv sessionP.SessionSrv) *Service {
	return &Service{
		UserRepo:       userRepo,
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
		SessionSrv:     sessionSrv,
	}
}

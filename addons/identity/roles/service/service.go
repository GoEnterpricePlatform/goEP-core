package service

import "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"

var _ port.RoleSrv = &Service{}

type Service struct {
	RoleRepo port.RoleRepo
}

func NewRoleSrv(roleRepo port.RoleRepo) *Service {
	return &Service{
		RoleRepo: roleRepo,
	}
}

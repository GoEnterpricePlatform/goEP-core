package service

import (
	"context"
	"fmt"

	permissionDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
	permissionPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/port"
	roleDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	rolePort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
)

type Service struct {
	PermissionRepo permissionPort.PermissionRepo
	RoleRepo       rolePort.RoleRepo
}

var _ permissionPort.PermissionSrv = (*Service)(nil)

func NewPermissionSrv(permissionRepo permissionPort.PermissionRepo, roleRepo rolePort.RoleRepo) *Service {
	return &Service{PermissionRepo: permissionRepo, RoleRepo: roleRepo}
}

func (s *Service) RegisterPermission(ctx context.Context, name permissionDomain.PermissionName, permissionType permissionDomain.PermissionType) error {
	if name == "" {
		return fmt.Errorf("permission name is required")
	}
	if permissionType == "" {
		return fmt.Errorf("permission type is required")
	}

	exists, err := s.PermissionRepo.Exists(ctx, name)
	if err != nil {
		return fmt.Errorf("check permission %q: %w", name, err)
	}
	if !exists {
		if err := s.PermissionRepo.Insert(ctx, permissionDomain.NewPermission(name, permissionType)); err != nil {
			return fmt.Errorf("insert permission %q: %w", name, err)
		}
	}

	permissions, err := s.PermissionRepo.FindByNames(ctx, []permissionDomain.PermissionName{name})
	if err != nil {
		return fmt.Errorf("load permission %q: %w", name, err)
	}
	if len(permissions) != 1 || permissions[0] == nil || permissions[0].ID == "" {
		return fmt.Errorf("permission %q was not found after registration", name)
	}
	if permissions[0].Type != permissionType {
		return fmt.Errorf("permission %q already exists with type %q, requested %q", name, permissions[0].Type, permissionType)
	}
	if err := s.RoleRepo.AssignPermissions(ctx, string(roleDomain.RoleSystemAdmin), []string{permissions[0].ID}); err != nil {
		return fmt.Errorf("assign permission %q to system administrator: %w", name, err)
	}
	return nil
}

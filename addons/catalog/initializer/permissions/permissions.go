package permissions

import (
	"context"

	permissionDomain "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
	permissionPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

func Register(services *modules.Services) error {
	permissionService, err := modules.RequireService[permissionPort.PermissionSrv](services)
	if err != nil {
		return err
	}
	return permissionService.RegisterPermission(context.Background(), permissionDomain.PermissionName("view.catalog"), permissionDomain.PermissionTypeView)
}

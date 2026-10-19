package config

import (
	adminPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/port"
	authPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	rolePort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	tokenPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	cookiePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/handler/cookie/port"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
)

// RegisterBackendModule connects Identity's typed constructor to the generic
// backend lifecycle and publishes the authentication contract it provides.
func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "identity",
		New:  NewIdentityModule,
	}); err != nil {
		return err
	}
	return modules.RegisterBackendInitializer(registry, "identity", func(runtime modules.BackendRuntime) (any, error) {
		module, err := NewIdentityModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			DB:         runtime.DB,
		})
		if err != nil {
			return nil, err
		}
		if err := modules.ProvideService[middlewarePort.Authenticator](runtime.Services, module.AuthApiMdw); err != nil {
			return nil, err
		}
		if err := modules.ProvideService[authPort.AuthSrv](runtime.Services, module.AuthSrv); err != nil {
			return nil, err
		}
		if err := modules.ProvideService[adminPort.AdminSrv](runtime.Services, module.AdminSrv); err != nil {
			return nil, err
		}
		if err := modules.ProvideService[tokenPort.TokenSrv](runtime.Services, module.TokenSrv); err != nil {
			return nil, err
		}
		if err := modules.ProvideService[cookiePort.CookieSrv](runtime.Services, module.CookieSrv); err != nil {
			return nil, err
		}
		if err := modules.ProvideService[rolePort.RoleSrv](runtime.Services, module.RoleSrv); err != nil {
			return nil, err
		}
		return module, nil
	})
}

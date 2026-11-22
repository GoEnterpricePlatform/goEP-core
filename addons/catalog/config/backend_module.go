package config

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/initializer/permissions"
	catalogPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	authPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	tokenPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	cookiePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/handler/cookie/port"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

// RegisterBackendModule connects Catalog to Identity and publishes the
// checkout resolver Billing uses to validate selected plan items.
func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "catalog",
		New:  NewCatalogModule,
	}); err != nil {
		return err
	}
	if err := modules.RegisterPermissionInitializer(registry, "catalog", permissions.Register); err != nil {
		return err
	}
	return modules.RegisterBackendInitializer(registry, "catalog", func(runtime modules.BackendRuntime) (any, error) {
		auth, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
		if err != nil {
			return nil, err
		}
		tokenSrv, err := modules.RequireService[tokenPort.TokenSrv](runtime.Services)
		if err != nil {
			return nil, err
		}
		authSrv, err := modules.RequireService[authPort.AuthSrv](runtime.Services)
		if err != nil {
			return nil, err
		}
		cookieSrv, err := modules.RequireService[cookiePort.CookieSrv](runtime.Services)
		if err != nil {
			return nil, err
		}
		module, err := NewCatalogModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			TemplateV1: runtime.TemplateV1,
			MdwSrvTmpl: middlewareServiceTmpl.NewMdwSrvTmpl(tokenSrv, authSrv, cookieSrv),
			DB:         runtime.DB,
			Deps:       ModuleDeps{AuthApiMdw: auth},
		})
		if err != nil {
			return nil, err
		}
		if err := modules.ProvideService[catalogPaddlePort.PaddleCheckoutCatalogSrv](runtime.Services, module.PaddleCheckoutCatalogSrv); err != nil {
			return nil, err
		}
		return module, nil
	})
}

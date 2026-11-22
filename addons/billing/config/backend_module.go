package config

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/initializer/permissions"
	paddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/port"
	catalogPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	authPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	tokenPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	cookiePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/handler/cookie/port"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "billing",
		New:  NewBillingModule,
	}); err != nil {
		return err
	}
	if err := modules.RegisterPermissionInitializer(registry, "billing", permissions.Register); err != nil {
		return err
	}

	return modules.RegisterBackendInitializer(registry, "billing", func(runtime modules.BackendRuntime) (any, error) {
		catalogCheckoutSrv, err := modules.RequireService[catalogPaddlePort.PaddleCheckoutCatalogSrv](runtime.Services)
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
		authMdw, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
		if err != nil {
			return nil, err
		}
		module, err := NewBillingModule(ModuleConfig{
			AppEnvs:            runtime.AppEnvs,
			AppClients:         runtime.AppClients,
			APIv1:              runtime.APIv1,
			AuthApiMdw:         authMdw,
			CatalogCheckoutSrv: catalogCheckoutSrv,
			TemplateV1:         runtime.TemplateV1,
			DB:                 runtime.DB,
			PGConn:             runtime.PGConn,
			MdwSrvTmpl:         middlewareServiceTmpl.NewMdwSrvTmpl(tokenSrv, authSrv, cookieSrv),
		})
		if err != nil {
			return nil, err
		}
		if err := modules.ProvideService[paddlePort.PaddleBillingSrv](runtime.Services, module.PaddleBillingSrv); err != nil {
			return nil, err
		}
		return module, nil
	})
}

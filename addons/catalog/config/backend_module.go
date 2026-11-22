package config

import (
	catalogPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
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
	return modules.RegisterBackendInitializer(registry, "catalog", func(runtime modules.BackendRuntime) (any, error) {
		auth, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
		if err != nil {
			return nil, err
		}
		module, err := NewCatalogModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
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

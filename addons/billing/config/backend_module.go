package config

import (
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
)

// RegisterBackendModule connects Billing to shared infrastructure and the
// authentication contract provided by Identity.
func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "billing",
		New:  NewBillingModule,
	}); err != nil {
		return err
	}
	return modules.RegisterBackendInitializer(registry, "billing", func(runtime modules.BackendRuntime) (any, error) {
		auth, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
		if err != nil {
			return nil, err
		}
		return NewBillingModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			DB:         runtime.DB,
			Deps:       ModuleDeps{AuthApiMdw: auth},
		})
	})
}

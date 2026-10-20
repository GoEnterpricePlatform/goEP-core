package config

import "github.com/GoEnterpricePlatform/goEP-core/internal/modules"

func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "invoices",
		New:  NewInvoicesModule,
	}); err != nil {
		return err
	}

	return modules.RegisterBackendInitializer(registry, "invoices", func(runtime modules.BackendRuntime) (any, error) {
		return NewInvoicesModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			DB:         runtime.DB,
		})
	})
}

package config

import (
	postPort "github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
)

// RegisterBackendModule connects Posts to shared infrastructure and the
// authentication contract provided by Identity.
func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "posts",
		New:  NewPostModule,
	}); err != nil {
		return err
	}
	return modules.RegisterBackendInitializer(registry, "posts", func(runtime modules.BackendRuntime) (any, error) {
		auth, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
		if err != nil {
			return nil, err
		}
		module, err := NewPostModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			DB:         runtime.DB,
			Deps:       ModuleDeps{AuthApiMdw: auth},
		})
		if err != nil {
			return nil, err
		}
		if err := modules.ProvideService[postPort.PostSrv](runtime.Services, module.PostService); err != nil {
			return nil, err
		}
		return module, nil
	})
}

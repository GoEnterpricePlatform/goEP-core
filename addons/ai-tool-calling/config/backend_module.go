package config

import (
	tcPort "github.com/GoEnterpricePlatform/goEP-core/addons/ai-tool-calling/port"
	postPort "github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

// RegisterBackendModule connects AI Tool Calling to shared infrastructure and
// the post service contract provided by Posts.
func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: "ai-tool-calling",
		New:  NewAiModule,
	}); err != nil {
		return err
	}
	return modules.RegisterBackendInitializer(registry, "ai-tool-calling", func(runtime modules.BackendRuntime) (any, error) {
		posts, err := modules.RequireService[postPort.PostSrv](runtime.Services)
		if err != nil {
			return nil, err
		}
		module, err := NewAiModule(ModuleConfig{
			AppEnvs:     runtime.AppEnvs,
			AppClients:  runtime.AppClients,
			APIv1:       runtime.APIv1,
			PostService: posts,
		})
		if err != nil {
			return nil, err
		}
		if err := modules.ProvideService[tcPort.ToolCallingSrv](runtime.Services, module.TCService); err != nil {
			return nil, err
		}
		return module, nil
	})
}

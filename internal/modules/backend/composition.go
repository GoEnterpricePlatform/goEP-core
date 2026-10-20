// Package backend owns the composition of the statically linked backend modules.
// It keeps constructor and dependency wiring details out of the HTTP bootstrap.
package backend

import (
	"fmt"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
)

// Config contains the infrastructure shared by backend modules and their route
// registrar. Each module still receives its own typed configuration.
type Config struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	AddonsDir  string
}

// Modules exposes module-provided contracts to other application layers.
type Modules struct {
	Services         *modules.Services
	Lifecycle        *lifecycle.Manager
	AvailableModules []modules.Manifest
}

// Initialize constructs backend modules in dependency order and registers
// their API routes on cfg.APIv1.
func Initialize(cfg Config) (*Modules, error) {
	registry, order, enabled, lifecycleManager, err := initializeModuleLifecycle(cfg)
	if err != nil {
		return nil, err
	}
	services := modules.NewServices()
	for _, manifest := range order {
		_, err := registry.InitializeBackendModule(manifest.Name, modules.BackendRuntime{
			AppEnvs:    cfg.AppEnvs,
			AppClients: cfg.AppClients,
			APIv1:      moduleRoutes(cfg.APIv1, enabled, manifest.Name),
			DB:         cfg.AppClients.DB,
			Services:   services,
		})
		if err != nil {
			return nil, fmt.Errorf("initialize backend module %q: %w", manifest.Name, err)
		}
	}

	return &Modules{
		Services:         services,
		Lifecycle:        lifecycleManager,
		AvailableModules: order,
	}, nil
}

func moduleRoutes(apiV1 *http.ServeMux, enabled map[string]bool, moduleName string) *http.ServeMux {
	if enabled[moduleName] {
		return apiV1
	}
	// Backend services are still constructed during this transition because
	// Web App consumes them directly. A private mux keeps disabled API routes
	// out of the public router until Web App joins the lifecycle manager.
	return http.NewServeMux()
}

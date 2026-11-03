// Package backend owns the composition of statically linked backend and web modules.
// It keeps constructor and dependency wiring details out of the HTTP bootstrap.
package backend

import (
	"fmt"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	//"github.com/jackc/pgx/v5"
)

// Config contains the infrastructure shared by backend modules and their route
// registrar. Each module still receives its own typed configuration.
type Config struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	WebMux     *http.ServeMux
	TemplateV1 *http.ServeMux
	AddonsDir  string
}

// Modules exposes module-provided contracts to other application layers.
type Modules struct {
	Services         *modules.Services
	Lifecycle        *lifecycle.Manager
	AvailableModules []modules.Manifest
}

// Initialize constructs backend modules in dependency order and then initializes
// web modules with the services they provide.
func Initialize(cfg Config) (*Modules, error) {
	registry, order, webManifests, enabled, lifecycleManager, err := initializeModuleLifecycle(cfg)
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
			PGConn:     cfg.AppClients.PGConn.DB,
			Services:   services,
		})
		if err != nil {
			return nil, fmt.Errorf("initialize backend module %q: %w", manifest.Name, err)
		}
	}
	if len(webManifests) > 0 {
		if cfg.WebMux == nil || cfg.TemplateV1 == nil {
			return nil, fmt.Errorf("web muxes are required to initialize web modules")
		}
		if err := registry.InitializeWebModules(webManifests, modules.WebRuntime{
			AppEnvs:          cfg.AppEnvs,
			Mux:              cfg.WebMux,
			TemplateV1:       cfg.TemplateV1,
			Services:         services,
			AvailableModules: order,
		}); err != nil {
			return nil, err
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
	// Backend services are constructed even when disabled because web modules
	// can depend on their service contracts. A private mux keeps their API routes isolated.
	return http.NewServeMux()
}

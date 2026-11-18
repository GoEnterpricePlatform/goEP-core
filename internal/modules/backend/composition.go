// Package backend owns the composition of statically linked backend and web modules.
// It keeps constructor and dependency wiring details out of the HTTP bootstrap.
package backend

import (
	"fmt"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend/lite"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend/managed"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Flavor selects how backend addons are composed at startup.
type Flavor string

const (
	// FlavorLite discovers and initializes all compiled addons in dependency order.
	// It does not persist install or enable state.
	FlavorLite Flavor = "lite"
	// FlavorManaged retains the install and enable lifecycle backed by the database.
	FlavorManaged Flavor = "managed"
)

// Config contains the infrastructure shared by backend modules and their route
// registrar. Each module still receives its own typed configuration.
type Config struct {
	// Flavor defaults to FlavorLite when empty.
	Flavor     Flavor
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
// web modules with the services they provide. An empty flavor selects FlavorLite.
func Initialize(cfg Config) (*Modules, error) {
	flavor := cfg.Flavor
	if flavor == "" {
		flavor = FlavorLite
	}

	var registry *modules.Registry
	var order, webManifests []modules.Manifest
	var enabled map[string]bool
	var lifecycleManager *lifecycle.Manager
	var err error
	compiledRegistry, err := NewCompiledRegistry()
	if err != nil {
		return nil, err
	}
	switch flavor {
	case FlavorLite:
		registry, order, webManifests, err = lite.Initialize(cfg.AppEnvs, cfg.AddonsDir, compiledRegistry)
		if err == nil {
			enabled = make(map[string]bool, len(order))
			for _, manifest := range order {
				enabled[manifest.Name] = true
			}
		}
	case FlavorManaged:
		registry, order, webManifests, enabled, lifecycleManager, err = managed.Initialize(cfg.AppEnvs, cfg.AppClients, cfg.AddonsDir, compiledRegistry)
	default:
		return nil, fmt.Errorf("unknown backend flavor %q; supported flavors are %q and %q", flavor, FlavorLite, FlavorManaged)
	}
	if err != nil {
		return nil, err
	}
	services := modules.NewServices()
	for _, manifest := range order {
		_, err := registry.InitializeBackendModule(manifest.Name, modules.BackendRuntime{
			AppEnvs:    cfg.AppEnvs,
			AppClients: cfg.AppClients,
			APIv1:      moduleRoutes(cfg.APIv1, enabled, manifest.Name),
			DB:         mongoDB(cfg.AppClients),
			PGConn:     postgresConnection(cfg.AppClients),
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

func mongoDB(clients *appConfig.AppClients) *mongo.Database {
	if clients == nil {
		return nil
	}
	return clients.DB
}

func postgresConnection(clients *appConfig.AppClients) *pgx.Conn {
	if clients == nil || clients.PGConn == nil {
		return nil
	}
	return clients.PGConn.DB
}

func moduleRoutes(apiV1 *http.ServeMux, enabled map[string]bool, moduleName string) *http.ServeMux {
	if enabled[moduleName] {
		return apiV1
	}
	// Backend services are constructed even when disabled because web modules
	// can depend on their service contracts. A private mux keeps their API routes isolated.
	return http.NewServeMux()
}

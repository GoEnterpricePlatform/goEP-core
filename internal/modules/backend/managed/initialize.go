// Package managed implements addon composition with persisted install/enable state.
package managed

import (
	"context"
	"fmt"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle/repository/mongo"
	postgres "github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle/repository/postgresql"
)

// Initialize discovers addons and prepares the database-backed lifecycle manager.
func Initialize(appEnvs *appConfig.AppEnvs, appClients *appConfig.AppClients, addonsDir string, registry *modules.Registry) (*modules.Registry, []modules.Manifest, []modules.Manifest, map[string]bool, *lifecycle.Manager, error) {
	if appEnvs == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("application environments are required to initialize module lifecycle")
	}
	switch appEnvs.DBProvider {
	case appConfig.DBMongo:
		if appClients == nil || appClients.DB == nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("MongoDB database is required to initialize module lifecycle")
		}
	case appConfig.DBPostgresql:
		if appClients == nil || appClients.PGConn == nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("PostgreSQL database is required to initialize module lifecycle")
		}
	}
	if addonsDir == "" {
		addonsDir = "addons"
	}
	allManifests, err := modules.NewDiscovery(addonsDir).DiscoverAll()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("discover modules: %w", err)
	}
	manifests, webManifests := splitModuleManifests(allManifests)
	if err := registry.ValidateBackend(manifests); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("validate backend module initializers: %w", err)
	}
	if err := registry.ValidateWeb(webManifests); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("validate web module initializers: %w", err)
	}
	manifests = filterManifestsByDatabase(manifests, string(appEnvs.DBProvider))
	if len(manifests) == 0 {
		return nil, nil, nil, nil, nil, fmt.Errorf("no backend modules support database provider %q", appEnvs.DBProvider)
	}
	registry, err = registry.Subset(manifests)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("select backend modules for database provider %q: %w", appEnvs.DBProvider, err)
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("resolve backend module dependencies: %w", err)
	}

	var storeRepo lifecycle.StateStore
	switch appEnvs.DBProvider {
	case appConfig.DBMongo:
		storeRepo = mongo.NewMongoStore(appClients.DB)
		if err := storeRepo.CreateIndexes(context.Background()); err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("initialize module lifecycle storage: %w", err)
		}
	case appConfig.DBPostgresql:
		storeRepo = postgres.NewPostgresStore(appClients.PGConn.Pool)
		if err := storeRepo.CreateIndexes(context.Background()); err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("initialize module lifecycle storage: %w", err)
		}
	}
	manager, err := lifecycle.NewManager(resolver, storeRepo)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	enabled, err := InitializeStates(context.Background(), manager, storeRepo, resolver)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return registry, resolver.Order(), webManifests, enabled, manager, nil
}

func splitModuleManifests(manifests []modules.Manifest) ([]modules.Manifest, []modules.Manifest) {
	backend := make([]modules.Manifest, 0, len(manifests))
	web := make([]modules.Manifest, 0)
	for _, manifest := range manifests {
		if manifest.Type == "web" {
			web = append(web, manifest)
			continue
		}
		backend = append(backend, manifest)
	}
	return backend, web
}

func filterManifestsByDatabase(manifests []modules.Manifest, provider string) []modules.Manifest {
	compatible := make([]modules.Manifest, 0, len(manifests))
	for _, manifest := range manifests {
		if manifest.SupportsDatabase(provider) {
			compatible = append(compatible, manifest)
		}
	}
	return compatible
}

// InitializeStates seeds a fresh managed installation and validates persisted state.
func InitializeStates(ctx context.Context, manager *lifecycle.Manager, store lifecycle.StateStore, resolver *modules.DependencyResolver) (map[string]bool, error) {
	states, err := store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read module lifecycle states: %w", err)
	}
	if len(states) == 0 {
		for _, manifest := range resolver.Order() {
			if err := manager.Install(ctx, manifest.Name); err != nil {
				return nil, fmt.Errorf("initialize module %q: %w", manifest.Name, err)
			}
		}
		for _, manifest := range resolver.Order() {
			if err := manager.Enable(ctx, manifest.Name); err != nil {
				return nil, fmt.Errorf("enable module %q on first startup: %w", manifest.Name, err)
			}
		}
	}
	statuses, err := manager.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list backend module statuses: %w", err)
	}
	enabledNames := make([]string, 0, len(statuses))
	enabled := make(map[string]bool, len(statuses))
	for _, status := range statuses {
		if status.Enabled && !status.Installed {
			return nil, fmt.Errorf("module %q is enabled but not installed", status.Manifest.Name)
		}
		if status.Enabled {
			enabled[status.Manifest.Name] = true
			enabledNames = append(enabledNames, status.Manifest.Name)
		}
	}
	for _, name := range enabledNames {
		if err := resolver.ValidateCanEnable(name, enabledNames); err != nil {
			return nil, fmt.Errorf("invalid enabled module state: %w", err)
		}
	}
	return enabled, nil
}

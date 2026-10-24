package backend

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle/repository/mongo"
	postgres "github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle/repository/postgresql"
)

func initializeModuleLifecycle(cfg Config) (*modules.Registry, []modules.Manifest, map[string]bool, *lifecycle.Manager, error) {
	switch cfg.AppEnvs.DBProvider {
	case config.DBMongo:
		if cfg.AppClients == nil || cfg.AppClients.DB == nil {
			return nil, nil, nil, nil, fmt.Errorf("MongoDB database is required to initialize module lifecycle")
		}
	case config.DBPostgresql:
		if cfg.AppClients == nil || cfg.AppClients.PGConn == nil {
			return nil, nil, nil, nil, fmt.Errorf("PostgreSQL database is required to initialize module lifecycle")
		}
	}

	addonsDir := cfg.AddonsDir
	if addonsDir == "" {
		addonsDir = "addons"
	}

	manifests, err := modules.NewDiscovery(addonsDir).Discover()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("discover backend modules: %w", err)
	}
	registry, err := NewCompiledRegistry()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := registry.ValidateBackend(manifests); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("validate backend module initializers: %w", err)
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("resolve backend module dependencies: %w", err)
	}

	var storeRepo lifecycle.StateStore
	switch cfg.AppEnvs.DBProvider {
	case config.DBMongo:
		storeRepo = mongo.NewMongoStore(cfg.AppClients.DB)
		if err := storeRepo.CreateIndexes(context.Background()); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("initialize module lifecycle storage: %w", err)
		}
	case config.DBPostgresql:
		if cfg.AppClients == nil || cfg.AppClients.PGConn == nil {
			return nil, nil, nil, nil, fmt.Errorf("PostgreSQL database is required to initialize module lifecycle")
		}
		storeRepo = postgres.NewPostgresStore(cfg.AppClients.PGConn.Pool)
	}

	manager, err := lifecycle.NewManager(resolver, storeRepo)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	enabled, err := initializeModuleStates(context.Background(), manager, storeRepo, resolver)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return registry, resolver.Order(), enabled, manager, nil
}

func initializeModuleStates(ctx context.Context, manager *lifecycle.Manager, store lifecycle.StateStore, resolver *modules.DependencyResolver) (map[string]bool, error) {
	states, err := store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read module lifecycle states: %w", err)
	}
	if len(states) == 0 {
		// Seed a fresh installation with the currently shipped backend modules
		// enabled so this migration preserves the existing server behavior.
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

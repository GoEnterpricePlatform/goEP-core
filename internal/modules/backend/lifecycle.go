package backend

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
)

func initializeModuleLifecycle(cfg Config) (*modules.Registry, map[string]bool, error) {
	if cfg.AppClients == nil || cfg.AppClients.DB == nil {
		return nil, nil, fmt.Errorf("MongoDB database is required to initialize module lifecycle")
	}
	addonsDir := cfg.AddonsDir
	if addonsDir == "" {
		addonsDir = "addons"
	}

	manifests, err := modules.NewDiscovery(addonsDir).Discover()
	if err != nil {
		return nil, nil, fmt.Errorf("discover backend modules: %w", err)
	}
	registry, err := NewCompiledRegistry()
	if err != nil {
		return nil, nil, err
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve backend module dependencies: %w", err)
	}

	store := lifecycle.NewMongoStore(cfg.AppClients.DB)
	if err := store.CreateIndexes(context.Background()); err != nil {
		return nil, nil, fmt.Errorf("initialize module lifecycle storage: %w", err)
	}
	manager, err := lifecycle.NewManager(resolver, store)
	if err != nil {
		return nil, nil, err
	}

	states, err := store.List(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("read module lifecycle states: %w", err)
	}
	if len(states) == 0 {
		// Seed a fresh installation with the currently shipped backend modules
		// enabled so this migration preserves the existing server behavior.
		for _, manifest := range resolver.Order() {
			if err := manager.Install(context.Background(), manifest.Name); err != nil {
				return nil, nil, fmt.Errorf("initialize module %q: %w", manifest.Name, err)
			}
		}
		for _, manifest := range resolver.Order() {
			if err := manager.Enable(context.Background(), manifest.Name); err != nil {
				return nil, nil, fmt.Errorf("enable module %q on first startup: %w", manifest.Name, err)
			}
		}
	}

	statuses, err := manager.List(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("list backend module statuses: %w", err)
	}
	enabledNames := make([]string, 0, len(statuses))
	enabled := make(map[string]bool, len(statuses))
	for _, status := range statuses {
		if status.Enabled && !status.Installed {
			return nil, nil, fmt.Errorf("module %q is enabled but not installed", status.Manifest.Name)
		}
		if status.Enabled {
			enabled[status.Manifest.Name] = true
			enabledNames = append(enabledNames, status.Manifest.Name)
		}
	}
	for _, name := range enabledNames {
		if err := resolver.ValidateCanEnable(name, enabledNames); err != nil {
			return nil, nil, fmt.Errorf("invalid enabled module state: %w", err)
		}
	}
	return registry, enabled, nil
}

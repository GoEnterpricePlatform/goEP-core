// Package lite implements the stateless addon composition flavor.
package lite

import (
	"fmt"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

// Initialize discovers compiled addons, validates their dependencies and
// returns them in startup order without reading or creating lifecycle state.
func Initialize(appEnvs *appConfig.AppEnvs, addonsDir string, registry *modules.Registry) (*modules.Registry, []modules.Manifest, []modules.Manifest, error) {
	if appEnvs == nil {
		return nil, nil, nil, fmt.Errorf("application environments are required to initialize backend modules")
	}
	if addonsDir == "" {
		addonsDir = "addons"
	}
	allManifests, err := modules.NewDiscovery(addonsDir).DiscoverAll()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("discover modules: %w", err)
	}
	backendManifests, webManifests := splitModuleManifests(allManifests)
	if err := registry.ValidateBackend(backendManifests); err != nil {
		return nil, nil, nil, fmt.Errorf("validate backend module initializers: %w", err)
	}
	if err := registry.ValidateWeb(webManifests); err != nil {
		return nil, nil, nil, fmt.Errorf("validate web module initializers: %w", err)
	}
	backendManifests = filterManifestsByDatabase(backendManifests, string(appEnvs.DBProvider))
	if len(backendManifests) == 0 {
		return nil, nil, nil, fmt.Errorf("no backend modules support database provider %q", appEnvs.DBProvider)
	}
	registry, err = registry.Subset(backendManifests)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("select backend modules for database provider %q: %w", appEnvs.DBProvider, err)
	}
	resolver, err := modules.NewDependencyResolver(backendManifests, registry)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve backend module dependencies: %w", err)
	}
	return registry, resolver.Order(), webManifests, nil
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

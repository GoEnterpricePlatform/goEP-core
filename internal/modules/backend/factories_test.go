package backend

import (
	"path/filepath"
	"testing"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

func TestCompiledBackendRegistryMatchesOfficialManifests(t *testing.T) {
	manifests, err := modules.NewDiscovery(filepath.Join("..", "..", "..", "addons")).Discover()
	if err != nil {
		t.Fatalf("Discover(repository addons) error = %v", err)
	}

	registry, err := NewCompiledRegistry()
	if err != nil {
		t.Fatalf("NewCompiledRegistry() error = %v", err)
	}
	if err := registry.ValidateBackend(manifests); err != nil {
		t.Fatalf("Registry.ValidateBackend(repository addons) error = %v", err)
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		t.Fatalf("NewDependencyResolver(repository addons) error = %v", err)
	}
	resolved := resolver.Order()
	position := make(map[string]int, len(resolved))
	for i, manifest := range resolved {
		position[manifest.Name] = i
	}
	for _, manifest := range manifests {
		for _, dependency := range manifest.Depends {
			dependencyPosition, exists := position[dependency]
			if !exists {
				t.Fatalf("resolved order is missing dependency %q of %q", dependency, manifest.Name)
			}
			if dependencyPosition >= position[manifest.Name] {
				t.Errorf("resolved order places %q before its dependency %q: %v", manifest.Name, dependency, moduleNames(resolved))
			}
		}
	}
}

func moduleNames(manifests []modules.Manifest) []string {
	names := make([]string, len(manifests))
	for i, manifest := range manifests {
		names[i] = manifest.Name
	}
	return names
}

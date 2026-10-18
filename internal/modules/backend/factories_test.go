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
	if err := registry.Validate(manifests); err != nil {
		t.Fatalf("Registry.Validate(repository addons) error = %v", err)
	}
	want := []string{"ai-tool-calling", "billing", "catalog", "identity", "posts"}
	got := registry.Names()
	if len(got) != len(want) {
		t.Fatalf("compiled registry names = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("compiled registry names = %#v, want %#v", got, want)
			break
		}
	}
}

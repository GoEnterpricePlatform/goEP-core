package modules

import (
	"strings"
	"testing"
)

func testRegistry(t *testing.T, names ...string) *Registry {
	t.Helper()
	registry := NewRegistry()
	for _, name := range names {
		if err := RegisterFactory(registry, Factory[string, string]{
			Name: name,
			New:  func(value string) (string, error) { return value, nil },
		}); err != nil {
			t.Fatal(err)
		}
	}
	return registry
}

func TestDependencyResolverReturnsDeterministicTopologicalOrder(t *testing.T) {
	manifests := []Manifest{
		{Name: "ai", Version: "1.0.0", Depends: []string{"posts"}},
		{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}},
		{Name: "catalog", Version: "1.0.0", Depends: []string{"identity"}},
		{Name: "identity", Version: "1.0.0"},
		{Name: "billing", Version: "1.0.0", Depends: []string{"identity"}},
	}
	resolver, err := NewDependencyResolver(manifests, testRegistry(t, "ai", "posts", "catalog", "identity", "billing"))
	if err != nil {
		t.Fatalf("NewDependencyResolver() error = %v", err)
	}

	want := []string{"identity", "billing", "catalog", "posts", "ai"}
	for run := 0; run < 3; run++ {
		order := resolver.Order()
		if len(order) != len(want) {
			t.Fatalf("Order() length = %d, want %d", len(order), len(want))
		}
		for i := range want {
			if order[i].Name != want[i] {
				t.Fatalf("Order() names = %v, want %v", manifestNames(order), want)
			}
		}
	}
}

func TestDependencyResolverRejectsMissingDependencies(t *testing.T) {
	manifests := []Manifest{{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}}}
	_, err := NewDependencyResolver(manifests, testRegistry(t, "posts"))
	if err == nil || !strings.Contains(err.Error(), `module "posts" depends on missing module "identity"`) {
		t.Fatalf("NewDependencyResolver() error = %v, want missing-dependency error", err)
	}
}

func TestDependencyResolverRejectsCycles(t *testing.T) {
	manifests := []Manifest{
		{Name: "alpha", Version: "1.0.0", Depends: []string{"beta"}},
		{Name: "beta", Version: "1.0.0", Depends: []string{"alpha"}},
	}
	_, err := NewDependencyResolver(manifests, testRegistry(t, "alpha", "beta"))
	if err == nil || !strings.Contains(err.Error(), "module dependency cycle detected") {
		t.Fatalf("NewDependencyResolver() error = %v, want cycle error", err)
	}
}

func TestDependencyResolverRequiresCompiledImplementations(t *testing.T) {
	manifests := []Manifest{{Name: "posts", Version: "1.0.0"}}
	_, err := NewDependencyResolver(manifests, NewRegistry())
	if err == nil || !strings.Contains(err.Error(), "has no compiled implementation") {
		t.Fatalf("NewDependencyResolver() error = %v, want missing-factory error", err)
	}
}

func TestDependencyResolverValidatesTransitiveEnabledDependencies(t *testing.T) {
	manifests := []Manifest{
		{Name: "identity", Version: "1.0.0"},
		{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}},
		{Name: "ai", Version: "1.0.0", Depends: []string{"posts"}},
	}
	resolver, err := NewDependencyResolver(manifests, testRegistry(t, "identity", "posts", "ai"))
	if err != nil {
		t.Fatalf("NewDependencyResolver() error = %v", err)
	}

	if err := resolver.ValidateCanEnable("ai", []string{"posts"}); err == nil || !strings.Contains(err.Error(), `dependency "identity" of "posts" is not enabled`) {
		t.Fatalf("ValidateCanEnable() error = %v, want transitive dependency error", err)
	}
	if err := resolver.ValidateCanEnable("ai", []string{"identity", "posts"}); err != nil {
		t.Fatalf("ValidateCanEnable(all dependencies) error = %v", err)
	}
	if err := resolver.ValidateCanEnable("unknown", nil); err == nil || !strings.Contains(err.Error(), "has no discovered manifest") {
		t.Fatalf("ValidateCanEnable(unknown) error = %v, want missing-manifest error", err)
	}
}

func manifestNames(manifests []Manifest) []string {
	names := make([]string, len(manifests))
	for i, manifest := range manifests {
		names[i] = manifest.Name
	}
	return names
}

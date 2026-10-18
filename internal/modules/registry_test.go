package modules

import (
	"strings"
	"testing"
)

func TestRegistryRegistersAndGetsTypedFactories(t *testing.T) {
	registry := NewRegistry()
	factory := Factory[string, int]{
		Name: "counter",
		New: func(value string) (int, error) {
			return len(value), nil
		},
	}
	if err := RegisterFactory(registry, factory); err != nil {
		t.Fatalf("RegisterFactory() error = %v", err)
	}

	registered, err := GetFactory[string, int](registry, "counter")
	if err != nil {
		t.Fatalf("GetFactory() error = %v", err)
	}
	got, err := registered.New("module")
	if err != nil {
		t.Fatalf("factory.New() error = %v", err)
	}
	if got != len("module") {
		t.Fatalf("factory.New() = %d, want %d", got, len("module"))
	}
	if names := registry.Names(); len(names) != 1 || names[0] != "counter" {
		t.Fatalf("Registry.Names() = %#v, want [counter]", names)
	}
}

func TestRegistryRejectsInvalidAndDuplicateFactories(t *testing.T) {
	registry := NewRegistry()
	tests := []struct {
		name    string
		factory Factory[string, int]
		want    string
	}{
		{name: "blank name", factory: Factory[string, int]{New: func(string) (int, error) { return 0, nil }}, want: "name is required"},
		{name: "nil constructor", factory: Factory[string, int]{Name: "missing", New: nil}, want: "has no constructor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterFactory(registry, tt.factory); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("RegisterFactory() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}

	factory := Factory[string, int]{Name: "counter", New: func(string) (int, error) { return 0, nil }}
	if err := RegisterFactory(registry, factory); err != nil {
		t.Fatalf("RegisterFactory() error = %v", err)
	}
	if err := RegisterFactory(registry, factory); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("RegisterFactory(duplicate) error = %v, want duplicate registration error", err)
	}
}

func TestGetFactoryRejectsMissingAndIncompatibleFactory(t *testing.T) {
	registry := NewRegistry()
	if err := RegisterFactory(registry, Factory[string, int]{Name: "counter", New: func(string) (int, error) { return 0, nil }}); err != nil {
		t.Fatal(err)
	}
	if _, err := GetFactory[string, int](registry, "unknown"); err == nil || !strings.Contains(err.Error(), "not registered") {
		t.Fatalf("GetFactory(unknown) error = %v, want not-registered error", err)
	}
	if _, err := GetFactory[int, string](registry, "counter"); err == nil || !strings.Contains(err.Error(), "incompatible types") {
		t.Fatalf("GetFactory(wrong types) error = %v, want incompatible-types error", err)
	}
}

func TestRegistryValidatesManifestsAgainstCompiledFactories(t *testing.T) {
	registry := NewRegistry()
	if err := RegisterFactory(registry, Factory[string, int]{Name: "identity", New: func(string) (int, error) { return 0, nil }}); err != nil {
		t.Fatal(err)
	}
	if err := RegisterFactory(registry, Factory[string, int]{Name: "posts", New: func(string) (int, error) { return 0, nil }}); err != nil {
		t.Fatal(err)
	}

	valid := []Manifest{
		{Name: "identity", Version: "1.0.0"},
		{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}},
	}
	if err := registry.Validate(valid); err != nil {
		t.Fatalf("Registry.Validate(valid) error = %v", err)
	}

	tests := []struct {
		name      string
		manifests []Manifest
		want      string
	}{
		{name: "uncompiled discovered module", manifests: append(valid, Manifest{Name: "catalog", Version: "1.0.0"}), want: `module "catalog" is discovered but has no compiled implementation`},
		{name: "compiled factory without manifest", manifests: valid[:1], want: `compiled module factory "posts" has no matching manifest`},
		{name: "duplicate manifest", manifests: append(valid, valid[0]), want: `duplicate module manifest "identity"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Validate(tt.manifests)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Registry.Validate() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

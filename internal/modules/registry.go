package modules

import (
	"fmt"
	"sort"
	"strings"
)

// Factory is a typed constructor for one compiled module implementation.
// Config and Implementation remain specific to each module.
type Factory[Config any, Implementation any] struct {
	Name string
	New  func(Config) (Implementation, error)
}

func (f Factory[Config, Implementation]) moduleName() string {
	return f.Name
}

type namedFactory interface {
	moduleName() string
}

// Registry contains the module implementations compiled into this executable.
type Registry struct {
	factories map[string]namedFactory
}

// NewRegistry creates an empty registry for statically compiled factories.
func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]namedFactory)}
}

// RegisterFactory adds a typed module factory to the registry.
func RegisterFactory[Config any, Implementation any](registry *Registry, factory Factory[Config, Implementation]) error {
	if registry == nil {
		return fmt.Errorf("module registry is nil")
	}
	name := strings.TrimSpace(factory.Name)
	if name == "" {
		return fmt.Errorf("module factory name is required")
	}
	if factory.New == nil {
		return fmt.Errorf("module factory %q has no constructor", name)
	}
	if _, exists := registry.factories[name]; exists {
		return fmt.Errorf("module factory %q is already registered", name)
	}
	factory.Name = name
	registry.factories[name] = factory
	return nil
}

// GetFactory retrieves a factory by name while preserving its module-specific
// configuration and implementation types.
func GetFactory[Config any, Implementation any](registry *Registry, name string) (Factory[Config, Implementation], error) {
	if registry == nil {
		return Factory[Config, Implementation]{}, fmt.Errorf("module registry is nil")
	}
	registered, exists := registry.factories[name]
	if !exists {
		return Factory[Config, Implementation]{}, fmt.Errorf("module factory %q is not registered", name)
	}
	factory, matchesType := registered.(Factory[Config, Implementation])
	if !matchesType {
		return Factory[Config, Implementation]{}, fmt.Errorf("module factory %q was requested with incompatible types", name)
	}
	return factory, nil
}

// Names returns registered implementation names in deterministic order.
func (r *Registry) Names() []string {
	if r == nil {
		return nil
	}
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Validate checks that every discovered module has a compiled factory and
// every registered factory has a corresponding manifest.
func (r *Registry) Validate(manifests []Manifest) error {
	if r == nil {
		return fmt.Errorf("module registry is nil")
	}

	manifestsByName := make(map[string]Manifest, len(manifests))
	ordered := append([]Manifest(nil), manifests...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })
	for _, manifest := range ordered {
		if err := manifest.Validate(); err != nil {
			return err
		}
		if _, exists := manifestsByName[manifest.Name]; exists {
			return fmt.Errorf("duplicate module manifest %q", manifest.Name)
		}
		manifestsByName[manifest.Name] = manifest
		if _, exists := r.factories[manifest.Name]; !exists {
			return fmt.Errorf("module %q is discovered but has no compiled implementation", manifest.Name)
		}
	}

	for _, name := range r.Names() {
		if _, exists := manifestsByName[name]; !exists {
			return fmt.Errorf("compiled module factory %q has no matching manifest", name)
		}
	}
	return nil
}

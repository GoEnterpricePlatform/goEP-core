package modules

import (
	"fmt"
	"sort"
)

// DependencyResolver validates module dependencies and stores a deterministic
// initialization order for the supplied manifests.
type DependencyResolver struct {
	manifests map[string]Manifest
	ordered   []Manifest
}

// NewDependencyResolver validates manifests, confirms that their dependencies
// are present and compiled, and computes a deterministic topological order.
func NewDependencyResolver(manifests []Manifest, registry *Registry) (*DependencyResolver, error) {
	byName := make(map[string]Manifest, len(manifests))
	for _, manifest := range manifests {
		if err := manifest.Validate(); err != nil {
			return nil, err
		}
		if _, exists := byName[manifest.Name]; exists {
			return nil, fmt.Errorf("duplicate module manifest %q", manifest.Name)
		}
		byName[manifest.Name] = manifest
	}

	for _, name := range sortedManifestNames(byName) {
		manifest := byName[name]
		for _, dependency := range manifest.Depends {
			if _, exists := byName[dependency]; !exists {
				return nil, fmt.Errorf("module %q depends on missing module %q", manifest.Name, dependency)
			}
		}
	}

	if err := registry.Validate(manifests); err != nil {
		return nil, err
	}
	ordered, err := topologicalOrder(byName)
	if err != nil {
		return nil, err
	}

	return &DependencyResolver{manifests: byName, ordered: ordered}, nil
}

// Order returns the resolved modules in dependency-first order.
func (r *DependencyResolver) Order() []Manifest {
	if r == nil {
		return nil
	}
	return append([]Manifest(nil), r.ordered...)
}

// ValidateCanEnable checks that every direct and transitive dependency of name
// is already enabled.
func (r *DependencyResolver) ValidateCanEnable(name string, enabledModules []string) error {
	if r == nil {
		return fmt.Errorf("module dependency resolver is nil")
	}
	if _, exists := r.manifests[name]; !exists {
		return fmt.Errorf("module %q has no discovered manifest", name)
	}

	enabled := make(map[string]struct{}, len(enabledModules))
	for _, module := range enabledModules {
		enabled[module] = struct{}{}
	}
	return r.validateDependenciesEnabled(name, name, enabled)
}

func (r *DependencyResolver) validateDependenciesEnabled(moduleName, rootName string, enabled map[string]struct{}) error {
	manifest := r.manifests[moduleName]
	dependencies := append([]string(nil), manifest.Depends...)
	sort.Strings(dependencies)
	for _, dependency := range dependencies {
		if _, ok := enabled[dependency]; !ok {
			if moduleName == rootName {
				return fmt.Errorf("module %q cannot be enabled: dependency %q is not enabled", rootName, dependency)
			}
			return fmt.Errorf("module %q cannot be enabled: dependency %q of %q is not enabled", rootName, dependency, moduleName)
		}
		if err := r.validateDependenciesEnabled(dependency, rootName, enabled); err != nil {
			return err
		}
	}
	return nil
}

func topologicalOrder(manifests map[string]Manifest) ([]Manifest, error) {
	indegree := make(map[string]int, len(manifests))
	dependents := make(map[string][]string, len(manifests))
	for name, manifest := range manifests {
		indegree[name] = len(manifest.Depends)
		for _, dependency := range manifest.Depends {
			dependents[dependency] = append(dependents[dependency], name)
		}
	}

	ready := make([]string, 0, len(manifests))
	for name, count := range indegree {
		if count == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)

	ordered := make([]Manifest, 0, len(manifests))
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		ordered = append(ordered, manifests[name])

		for _, dependent := range dependents[name] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.Strings(ready)
	}

	if len(ordered) != len(manifests) {
		cycle := make([]string, 0)
		for name, count := range indegree {
			if count > 0 {
				cycle = append(cycle, name)
			}
		}
		sort.Strings(cycle)
		return nil, fmt.Errorf("module dependency cycle detected among %v", cycle)
	}
	return ordered, nil
}

func sortedManifestNames(manifests map[string]Manifest) []string {
	names := make([]string, 0, len(manifests))
	for name := range manifests {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

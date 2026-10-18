package lifecycle

import (
	"context"
	"fmt"
	"sync"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

// Manager coordinates installation and activation while leaving module data
// untouched when a module is disabled.
type Manager struct {
	resolver *modules.DependencyResolver
	store    StateStore
	mu       sync.Mutex
}

// NewManager creates a lifecycle manager for discovered, compiled modules.
func NewManager(resolver *modules.DependencyResolver, store StateStore) (*Manager, error) {
	if resolver == nil {
		return nil, fmt.Errorf("module dependency resolver is required")
	}
	if store == nil {
		return nil, fmt.Errorf("module state store is required")
	}
	return &Manager{resolver: resolver, store: store}, nil
}

// List returns discovered modules and their current installation state.
func (m *Manager) List(ctx context.Context) ([]ModuleStatus, error) {
	states, err := m.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list module states: %w", err)
	}
	stateByName := make(map[string]ModuleState, len(states))
	for _, state := range states {
		if state.Enabled && !state.Installed {
			return nil, fmt.Errorf("module state %q is enabled but not installed", state.Name)
		}
		stateByName[state.Name] = state
	}

	ordered := m.resolver.Order()
	statuses := make([]ModuleStatus, 0, len(ordered))
	for _, manifest := range ordered {
		state := stateByName[manifest.Name]
		statuses = append(statuses, ModuleStatus{
			Manifest:  manifest,
			Installed: state.Installed,
			Enabled:   state.Enabled,
		})
	}
	return statuses, nil
}

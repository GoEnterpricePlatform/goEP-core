package lifecycle

import (
	"context"
	"fmt"
	"time"
)

// Enable activates an installed module after confirming all direct and
// transitive dependencies are installed and enabled.
func (m *Manager) Enable(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.resolver.Manifest(name); !exists {
		return fmt.Errorf("%w: %q", ErrModuleNotAvailable, name)
	}
	state, found, err := m.store.Find(ctx, name)
	if err != nil {
		return fmt.Errorf("read module state %q: %w", name, err)
	}
	if !found || !state.Installed {
		return fmt.Errorf("%w: %q", ErrModuleNotInstalled, name)
	}
	if state.Enabled {
		return nil
	}

	states, err := m.store.List(ctx)
	if err != nil {
		return fmt.Errorf("list module states: %w", err)
	}
	stateByName := make(map[string]ModuleState, len(states))
	enabled := make([]string, 0, len(states))
	for _, current := range states {
		stateByName[current.Name] = current
		if current.Installed && current.Enabled {
			enabled = append(enabled, current.Name)
		}
	}
	for _, dependency := range m.resolver.RequiredDependencies(name) {
		dependencyState, found := stateByName[dependency]
		if !found || !dependencyState.Installed {
			return fmt.Errorf("module %q cannot be enabled: dependency %q is not installed", name, dependency)
		}
		if !dependencyState.Enabled {
			return fmt.Errorf("module %q cannot be enabled: dependency %q is not enabled", name, dependency)
		}
	}
	if err := m.resolver.ValidateCanEnable(name, enabled); err != nil {
		return err
	}

	state.Enabled = true
	state.UpdatedAt = time.Now().UTC()
	if err := m.store.Save(ctx, state); err != nil {
		return fmt.Errorf("enable module %q: %w", name, err)
	}
	return nil
}

// Disable deactivates an installed module without deleting module data or its
// installation record. It refuses to break an enabled dependent module.
func (m *Manager) Disable(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.resolver.Manifest(name); !exists {
		return fmt.Errorf("%w: %q", ErrModuleNotAvailable, name)
	}
	state, found, err := m.store.Find(ctx, name)
	if err != nil {
		return fmt.Errorf("read module state %q: %w", name, err)
	}
	if !found || !state.Installed || !state.Enabled {
		return nil
	}
	states, err := m.store.List(ctx)
	if err != nil {
		return fmt.Errorf("list module states: %w", err)
	}
	enabled := make([]string, 0, len(states))
	for _, current := range states {
		if current.Installed && current.Enabled {
			enabled = append(enabled, current.Name)
		}
	}
	if err := m.resolver.ValidateCanDisable(name, enabled); err != nil {
		return err
	}

	state.Enabled = false
	state.UpdatedAt = time.Now().UTC()
	if err := m.store.Save(ctx, state); err != nil {
		return fmt.Errorf("disable module %q: %w", name, err)
	}
	return nil
}

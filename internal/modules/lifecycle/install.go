package lifecycle

import (
	"context"
	"fmt"
	"time"
)

// Install marks a discovered module installed. Reinstalling is idempotent;
// when the compiled manifest version changes, it updates the stored version
// while preserving the enabled flag.
func (m *Manager) Install(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	manifest, exists := m.resolver.Manifest(name)
	if !exists {
		return fmt.Errorf("%w: %q", ErrModuleNotAvailable, name)
	}
	state, found, err := m.store.Find(ctx, name)
	if err != nil {
		return fmt.Errorf("read module state %q: %w", name, err)
	}
	now := time.Now().UTC()
	if found && state.Installed {
		if state.Version == manifest.Version {
			return nil
		}
		if state.Enabled {
			states, err := m.store.List(ctx)
			if err != nil {
				return fmt.Errorf("list module states before updating %q: %w", name, err)
			}
			enabled := make([]string, 0, len(states))
			for _, current := range states {
				if current.Installed && current.Enabled {
					enabled = append(enabled, current.Name)
				}
			}
			if err := m.resolver.ValidateCanEnable(name, enabled); err != nil {
				return fmt.Errorf("cannot update enabled module %q until its dependencies are enabled: %w", name, err)
			}
		}
		state.Version = manifest.Version
		state.UpdatedAt = now
		if err := m.store.Save(ctx, state); err != nil {
			return fmt.Errorf("update module %q to version %q: %w", name, manifest.Version, err)
		}
		return nil
	}

	if !found || state.InstalledAt.IsZero() {
		state.InstalledAt = now
	}
	state.Name = name
	state.Version = manifest.Version
	state.Installed = true
	state.Enabled = false
	state.UpdatedAt = now
	if err := m.store.Save(ctx, state); err != nil {
		return fmt.Errorf("install module %q: %w", name, err)
	}
	return nil
}

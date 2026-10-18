package lifecycle

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

type memoryStateStore struct {
	states map[string]ModuleState
}

func newMemoryStateStore() *memoryStateStore {
	return &memoryStateStore{states: make(map[string]ModuleState)}
}

func (s *memoryStateStore) Find(_ context.Context, name string) (ModuleState, bool, error) {
	state, exists := s.states[name]
	return state, exists, nil
}

func (s *memoryStateStore) List(_ context.Context) ([]ModuleState, error) {
	states := make([]ModuleState, 0, len(s.states))
	for _, state := range s.states {
		states = append(states, state)
	}
	return states, nil
}

func (s *memoryStateStore) Save(_ context.Context, state ModuleState) error {
	s.states[state.Name] = state
	return nil
}

func lifecycleManager(t *testing.T) (*Manager, *memoryStateStore) {
	t.Helper()
	manifests := []modules.Manifest{
		{Name: "identity", Version: "1.0.0"},
		{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}},
		{Name: "ai", Version: "2.0.0", Depends: []string{"posts"}},
	}
	registry := modules.NewRegistry()
	for _, manifest := range manifests {
		name := manifest.Name
		if err := modules.RegisterFactory(registry, modules.Factory[string, string]{Name: name, New: func(value string) (string, error) { return value, nil }}); err != nil {
			t.Fatal(err)
		}
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		t.Fatal(err)
	}
	store := newMemoryStateStore()
	manager, err := NewManager(resolver, store)
	if err != nil {
		t.Fatal(err)
	}
	return manager, store
}

func TestManagerListsAvailableInstalledAndEnabledStates(t *testing.T) {
	manager, _ := lifecycleManager(t)
	ctx := context.Background()
	if err := manager.Install(ctx, "identity"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Enable(ctx, "identity"); err != nil {
		t.Fatal(err)
	}

	statuses, err := manager.List(ctx)
	if err != nil {
		t.Fatalf("Manager.List() error = %v", err)
	}
	want := map[string]ModuleStatus{
		"identity": {Manifest: modules.Manifest{Name: "identity", Version: "1.0.0"}, Installed: true, Enabled: true},
		"posts":    {Manifest: modules.Manifest{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}}},
		"ai":       {Manifest: modules.Manifest{Name: "ai", Version: "2.0.0", Depends: []string{"posts"}}},
	}
	if len(statuses) != len(want) {
		t.Fatalf("Manager.List() returned %d statuses, want %d", len(statuses), len(want))
	}
	for _, status := range statuses {
		expected, exists := want[status.Manifest.Name]
		if !exists || status.Installed != expected.Installed || status.Enabled != expected.Enabled {
			t.Errorf("unexpected module status: %#v", status)
		}
	}
}

func TestManagerInstallIsIdempotentAndUpdatesManifestVersion(t *testing.T) {
	manager, store := lifecycleManager(t)
	ctx := context.Background()
	if err := manager.Install(ctx, "ai"); err != nil {
		t.Fatal(err)
	}
	first := store.states["ai"]
	if !first.Installed || first.Enabled || first.Version != "2.0.0" {
		t.Fatalf("installed state = %#v, unexpected initial state", first)
	}
	if err := manager.Install(ctx, "ai"); err != nil {
		t.Fatalf("repeated Install() error = %v", err)
	}
	if got := store.states["ai"].UpdatedAt; !got.Equal(first.UpdatedAt) {
		t.Fatalf("repeated Install() changed UpdatedAt from %v to %v", first.UpdatedAt, got)
	}

	for _, name := range []string{"identity", "posts"} {
		if err := manager.Install(ctx, name); err != nil {
			t.Fatal(err)
		}
		if err := manager.Enable(ctx, name); err != nil {
			t.Fatal(err)
		}
	}
	store.states["ai"] = ModuleState{Name: "ai", Version: "1.0.0", Installed: true, Enabled: true}
	if err := manager.Install(ctx, "ai"); err != nil {
		t.Fatal(err)
	}
	updated := store.states["ai"]
	if updated.Version != "2.0.0" || !updated.Enabled {
		t.Fatalf("updated state = %#v, want current version preserving enabled state", updated)
	}
}

func TestManagerEnableRequiresInstallationAndEnabledDependencies(t *testing.T) {
	manager, _ := lifecycleManager(t)
	ctx := context.Background()
	if err := manager.Install(ctx, "posts"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Enable(ctx, "posts"); err == nil || !strings.Contains(err.Error(), `dependency "identity" is not installed`) {
		t.Fatalf("Enable(posts) error = %v, want dependency installation error", err)
	}
	if err := manager.Install(ctx, "identity"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Enable(ctx, "posts"); err == nil || !strings.Contains(err.Error(), `dependency "identity" is not enabled`) {
		t.Fatalf("Enable(posts) error = %v, want dependency activation error", err)
	}
	if err := manager.Enable(ctx, "identity"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Enable(ctx, "posts"); err != nil {
		t.Fatalf("Enable(posts) with identity enabled error = %v", err)
	}
	if err := manager.Enable(ctx, "not-available"); !errors.Is(err, ErrModuleNotAvailable) {
		t.Fatalf("Enable(unknown) error = %v, want unavailable error", err)
	}
}

func TestManagerDisablePreservesInstallationAndProtectsEnabledDependents(t *testing.T) {
	manager, store := lifecycleManager(t)
	ctx := context.Background()
	for _, name := range []string{"identity", "posts", "ai"} {
		if err := manager.Install(ctx, name); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"identity", "posts", "ai"} {
		if err := manager.Enable(ctx, name); err != nil {
			t.Fatal(err)
		}
	}

	if err := manager.Disable(ctx, "posts"); err == nil {
		t.Fatal("Disable(posts) succeeded while ai was enabled")
	}
	if err := manager.Disable(ctx, "ai"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Disable(ctx, "posts"); err != nil {
		t.Fatal(err)
	}
	state, found, err := store.Find(ctx, "posts")
	if err != nil {
		t.Fatal(err)
	}
	if !found || !state.Installed || state.Enabled {
		t.Fatalf("disabled posts state = %#v, found=%t; want installed and disabled", state, found)
	}
	if err := manager.Disable(ctx, "not-available"); !errors.Is(err, ErrModuleNotAvailable) {
		t.Fatalf("Disable(unknown) error = %v, want unavailable error", err)
	}
}

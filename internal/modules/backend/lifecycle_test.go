package backend

import (
	"context"
	"testing"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
)

type testLifecycleStore struct {
	states map[string]lifecycle.ModuleState
}

func newTestLifecycleStore() *testLifecycleStore {
	return &testLifecycleStore{states: make(map[string]lifecycle.ModuleState)}
}

func (s *testLifecycleStore) Find(_ context.Context, name string) (lifecycle.ModuleState, bool, error) {
	state, exists := s.states[name]
	return state, exists, nil
}

func (s *testLifecycleStore) List(_ context.Context) ([]lifecycle.ModuleState, error) {
	states := make([]lifecycle.ModuleState, 0, len(s.states))
	for _, state := range s.states {
		states = append(states, state)
	}
	return states, nil
}

func (s *testLifecycleStore) Save(_ context.Context, state lifecycle.ModuleState) error {
	s.states[state.Name] = state
	return nil
}

func testLifecycleManager(t *testing.T, store lifecycle.StateStore) (*lifecycle.Manager, *modules.DependencyResolver) {
	t.Helper()
	manifests := []modules.Manifest{
		{Name: "identity", Version: "1.0.0"},
		{Name: "posts", Version: "1.0.0", Depends: []string{"identity"}},
		{Name: "ai-tool-calling", Version: "1.0.0", Depends: []string{"posts"}},
	}
	registry := modules.NewRegistry()
	for _, manifest := range manifests {
		if err := modules.RegisterFactory(registry, modules.Factory[string, string]{
			Name: manifest.Name,
			New:  func(value string) (string, error) { return value, nil },
		}); err != nil {
			t.Fatal(err)
		}
	}
	resolver, err := modules.NewDependencyResolver(manifests, registry)
	if err != nil {
		t.Fatal(err)
	}
	manager, err := lifecycle.NewManager(resolver, store)
	if err != nil {
		t.Fatal(err)
	}
	return manager, resolver
}

func TestInitializeModuleStatesEnablesModulesOnEmptyStore(t *testing.T) {
	store := newTestLifecycleStore()
	manager, resolver := testLifecycleManager(t, store)

	enabled, err := initializeModuleStates(context.Background(), manager, store, resolver)
	if err != nil {
		t.Fatalf("initializeModuleStates() error = %v", err)
	}
	if len(enabled) != 3 {
		t.Fatalf("enabled modules = %#v, want all three modules", enabled)
	}
	for _, name := range []string{"identity", "posts", "ai-tool-calling"} {
		if !enabled[name] {
			t.Errorf("module %q is disabled after initial setup", name)
		}
		if state := store.states[name]; !state.Installed || !state.Enabled {
			t.Errorf("initial state for %q = %#v, want installed and enabled", name, state)
		}
	}
}

func TestInitializeModuleStatesHonorsPersistedStateAndRejectsBrokenDependencies(t *testing.T) {
	store := newTestLifecycleStore()
	store.states["posts"] = lifecycle.ModuleState{
		Name:      "posts",
		Version:   "1.0.0",
		Installed: true,
		Enabled:   true,
	}
	manager, resolver := testLifecycleManager(t, store)

	if _, err := initializeModuleStates(context.Background(), manager, store, resolver); err == nil {
		t.Fatal("initializeModuleStates() succeeded with posts enabled and its identity dependency disabled")
	}
	if _, exists := store.states["identity"]; exists {
		t.Fatal("initializeModuleStates() seeded defaults despite existing persisted state")
	}
}

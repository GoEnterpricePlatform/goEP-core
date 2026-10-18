package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoveryReturnsValidatedManifestsSortedByName(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, "zeta", `{"name":"zeta","version":"1.0.0"}`)
	writeManifest(t, root, "alpha", `{"name":"alpha","version":"1.0.0"}`)
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("not an addon directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Discovery treats Go files as ordinary files and never tries to execute them.
	if err := os.WriteFile(filepath.Join(root, "zeta", "module.go"), []byte("this is not valid Go"), 0o600); err != nil {
		t.Fatal(err)
	}

	manifests, err := NewDiscovery(root).Discover()
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(manifests) != 2 {
		t.Fatalf("Discover() returned %d manifests, want 2", len(manifests))
	}
	if manifests[0].Name != "alpha" || manifests[1].Name != "zeta" {
		t.Fatalf("Discover() names = [%q, %q], want [alpha, zeta]", manifests[0].Name, manifests[1].Name)
	}
}

func TestDiscoveryFindsRepositoryAddons(t *testing.T) {
	manifests, err := NewDiscovery(filepath.Join("..", "..", "addons")).Discover()
	if err != nil {
		t.Fatalf("Discover(repository addons) error = %v", err)
	}
	want := []string{"ai-tool-calling", "billing", "catalog", "identity", "posts"}
	if len(manifests) != len(want) {
		t.Fatalf("Discover(repository addons) returned %d modules, want %d", len(manifests), len(want))
	}
	for i := range want {
		if manifests[i].Name != want[i] {
			t.Errorf("Discover(repository addons)[%d].Name = %q, want %q", i, manifests[i].Name, want[i])
		}
	}
}

func TestDiscoveryRejectsDirectoryWithoutManifest(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "unfinished"), 0o700); err != nil {
		t.Fatal(err)
	}

	_, err := NewDiscovery(root).Discover()
	if err == nil || !strings.Contains(err.Error(), "unfinished") || !strings.Contains(err.Error(), manifestFilename) {
		t.Fatalf("Discover() error = %v, want missing-manifest error for unfinished addon", err)
	}
}

func TestDiscoveryRejectsInvalidManifest(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, "invalid", `{"name":"invalid"}`)

	_, err := NewDiscovery(root).Discover()
	if err == nil || !strings.Contains(err.Error(), "version is required") {
		t.Fatalf("Discover() error = %v, want missing-version error", err)
	}
}

func TestDiscoveryRejectsDuplicateModuleNames(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, "one", `{"name":"shared","version":"1.0.0"}`)
	writeManifest(t, root, "two", `{"name":"shared","version":"1.0.0"}`)

	_, err := NewDiscovery(root).Discover()
	if err == nil || !strings.Contains(err.Error(), `duplicate module name "shared"`) {
		t.Fatalf("Discover() error = %v, want duplicate-name error", err)
	}
}

func writeManifest(t *testing.T, root, dirName, manifestJSON string) {
	t.Helper()
	addonDir := filepath.Join(root, dirName)
	if err := os.Mkdir(addonDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(addonDir, manifestFilename), []byte(manifestJSON), 0o600); err != nil {
		t.Fatal(err)
	}
}

package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseManifest(t *testing.T) {
	manifest, err := ParseManifest([]byte(`{
		"name": "posts",
		"version": "1.0.0",
		"summary": "Post management",
		"depends": ["identity"]
	}`))
	if err != nil {
		t.Fatalf("ParseManifest() error = %v", err)
	}
	if manifest.Name != "posts" || manifest.Version != "1.0.0" || manifest.Summary != "Post management" {
		t.Fatalf("ParseManifest() = %#v, unexpected metadata", manifest)
	}
	if len(manifest.Depends) != 1 || manifest.Depends[0] != "identity" {
		t.Fatalf("ParseManifest() depends = %#v, want [identity]", manifest.Depends)
	}
}

func TestParseManifestRequiresNameAndVersion(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{name: "missing name", json: `{"version":"1.0.0"}`, want: "name is required"},
		{name: "blank name", json: `{"name":"  ","version":"1.0.0"}`, want: "name is required"},
		{name: "missing version", json: `{"name":"posts"}`, want: "version is required"},
		{name: "blank version", json: `{"name":"posts","version":" "}`, want: "version is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseManifest([]byte(tt.json))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ParseManifest() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestParseManifestRejectsInvalidJSONAndInconsistentDependencies(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{name: "malformed JSON", json: `{"name":`, want: "decode module manifest"},
		{name: "unknown field", json: `{"name":"posts","version":"1.0.0","depend":"identity"}`, want: "unknown field"},
		{name: "trailing JSON", json: `{"name":"posts","version":"1.0.0"} {}`, want: "multiple JSON values"},
		{name: "empty dependency", json: `{"name":"posts","version":"1.0.0","depends":[""]}`, want: "depends entries must not be empty"},
		{name: "self dependency", json: `{"name":"posts","version":"1.0.0","depends":["posts"]}`, want: "cannot depend on itself"},
		{name: "duplicate dependency", json: `{"name":"posts","version":"1.0.0","depends":["identity","identity"]}`, want: "duplicate dependency"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseManifest([]byte(tt.json))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ParseManifest() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestAddonManifests(t *testing.T) {
	tests := []struct {
		name    string
		depends []string
	}{
		{name: "identity"},
		{name: "posts", depends: []string{"identity"}},
		{name: "billing", depends: []string{"identity"}},
		{name: "catalog", depends: []string{"identity"}},
		{name: "ai-tool-calling", depends: []string{"posts"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "addons", tt.name, "manifest.json")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", path, err)
			}

			manifest, err := ParseManifest(data)
			if err != nil {
				t.Fatalf("ParseManifest(%q) error = %v", path, err)
			}
			if manifest.Name != tt.name {
				t.Errorf("manifest name = %q, want %q", manifest.Name, tt.name)
			}
			if len(manifest.Depends) != len(tt.depends) {
				t.Fatalf("manifest depends = %#v, want %#v", manifest.Depends, tt.depends)
			}
			for i := range tt.depends {
				if manifest.Depends[i] != tt.depends[i] {
					t.Errorf("manifest depends = %#v, want %#v", manifest.Depends, tt.depends)
					break
				}
			}
		})
	}
}

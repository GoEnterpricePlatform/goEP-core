package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
)

// Manifest describes a GoEP module independently from its Go implementation.
// An omitted Type is treated as "backend" for existing addon manifests.
type Manifest struct {
	Type              string   `json:"type,omitempty"`
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	Summary           string   `json:"summary"`
	Icon              string   `json:"icon,omitempty"`
	Depends           []string `json:"depends"`
	DatabaseProviders []string `json:"database_providers,omitempty"`
}

// SupportsDatabase reports whether the module declares support for a database
// provider. An empty list preserves compatibility with manifests created
// before database provider support was declared.
func (m Manifest) SupportsDatabase(provider string) bool {
	if len(m.DatabaseProviders) == 0 {
		return true
	}
	for _, supported := range m.DatabaseProviders {
		if supported == provider {
			return true
		}
	}
	return false
}

// ParseManifest decodes and validates a manifest encoded as JSON.
func ParseManifest(data []byte) (Manifest, error) {
	return ReadManifest(bytes.NewReader(data))
}

// ReadManifest decodes and validates one manifest from r.
func ReadManifest(r io.Reader) (Manifest, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode module manifest: %w", err)
	}
	if err := ensureJSONEnd(decoder); err != nil {
		return Manifest{}, err
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// Validate checks required metadata and basic dependency-list consistency.
func (m Manifest) Validate() error {
	if m.Type != "" && m.Type != "backend" && m.Type != "web" {
		return fmt.Errorf("module manifest %q: type must be \"backend\" or \"web\"", m.Name)
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("module manifest: name is required")
	}
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("module manifest %q: version is required", m.Name)
	}
	if m.Icon != "" {
		cleanIcon := path.Clean(m.Icon)
		if strings.Contains(m.Icon, `\`) || cleanIcon != m.Icon || !strings.HasPrefix(cleanIcon, "static/public/") {
			return fmt.Errorf("module manifest %q: icon must be a clean path inside static/public/", m.Name)
		}
	}

	seen := make(map[string]struct{}, len(m.Depends))
	for _, dependency := range m.Depends {
		dependency = strings.TrimSpace(dependency)
		if dependency == "" {
			return fmt.Errorf("module manifest %q: depends entries must not be empty", m.Name)
		}
		if dependency == m.Name {
			return fmt.Errorf("module manifest %q: module cannot depend on itself", m.Name)
		}
		if _, exists := seen[dependency]; exists {
			return fmt.Errorf("module manifest %q: duplicate dependency %q", m.Name, dependency)
		}
		seen[dependency] = struct{}{}
	}

	seenProviders := make(map[string]struct{}, len(m.DatabaseProviders))
	for _, provider := range m.DatabaseProviders {
		provider = strings.TrimSpace(provider)
		if provider == "" {
			return fmt.Errorf("module manifest %q: database_providers entries must not be empty", m.Name)
		}
		if _, exists := seenProviders[provider]; exists {
			return fmt.Errorf("module manifest %q: duplicate database provider %q", m.Name, provider)
		}
		seenProviders[provider] = struct{}{}
	}
	return nil
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode module manifest: multiple JSON values are not allowed")
		}
		return fmt.Errorf("decode module manifest: %w", err)
	}
	return nil
}

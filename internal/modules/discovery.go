package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const manifestFilename = "manifest.json"

// Discovery reads backend module metadata from the immediate child directories of root.
type Discovery struct {
	root string
}

// NewDiscovery creates a discovery service for the given addons directory.
func NewDiscovery(root string) *Discovery {
	return &Discovery{root: root}
}

// Discover validates child manifests and returns backend modules. A directory
// without manifest.json is treated as an incomplete addon and causes an error.
// Web modules are validated but omitted. Discovery never loads or executes Go source files.
func (d *Discovery) Discover() ([]Manifest, error) {
	manifests, err := d.DiscoverAll()
	if err != nil {
		return nil, err
	}
	backend := make([]Manifest, 0, len(manifests))
	for _, manifest := range manifests {
		if manifest.Type != "web" {
			backend = append(backend, manifest)
		}
	}
	return backend, nil
}

// DiscoverAll validates and returns backend and web module manifests sorted by name.
func (d *Discovery) DiscoverAll() ([]Manifest, error) {
	rootInfo, err := os.Stat(d.root)
	if err != nil {
		return nil, fmt.Errorf("inspect addons directory %q: %w", d.root, err)
	}
	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("addons path %q is not a directory", d.root)
	}

	entries, err := os.ReadDir(d.root)
	if err != nil {
		return nil, fmt.Errorf("read addons directory %q: %w", d.root, err)
	}

	manifests := make([]Manifest, 0, len(entries))
	manifestDirs := make(map[string]string, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		addonDir := filepath.Join(d.root, entry.Name())
		manifestPath := filepath.Join(addonDir, manifestFilename)
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("addon directory %q is missing %s", addonDir, manifestFilename)
			}
			return nil, fmt.Errorf("read manifest %q: %w", manifestPath, err)
		}

		manifest, err := ParseManifest(data)
		if err != nil {
			return nil, fmt.Errorf("invalid manifest %q: %w", manifestPath, err)
		}
		if previousDir, exists := manifestDirs[manifest.Name]; exists {
			return nil, fmt.Errorf("duplicate module name %q in %q and %q", manifest.Name, previousDir, addonDir)
		}
		manifestDirs[manifest.Name] = addonDir
		if manifest.Icon != "" {
			iconPath := filepath.Join(addonDir, filepath.FromSlash(manifest.Icon))
			iconInfo, err := os.Stat(iconPath)
			if err != nil {
				return nil, fmt.Errorf("module %q icon %q: %w", manifest.Name, manifest.Icon, err)
			}
			if !iconInfo.Mode().IsRegular() {
				return nil, fmt.Errorf("module %q icon %q is not a regular file", manifest.Name, manifest.Icon)
			}
			addonReal, err := filepath.EvalSymlinks(addonDir)
			if err != nil {
				return nil, fmt.Errorf("resolve addon directory %q: %w", addonDir, err)
			}
			iconReal, err := filepath.EvalSymlinks(iconPath)
			if err != nil {
				return nil, fmt.Errorf("resolve module %q icon: %w", manifest.Name, err)
			}
			relativeIcon, err := filepath.Rel(addonReal, iconReal)
			if err != nil || relativeIcon == ".." || strings.HasPrefix(relativeIcon, ".."+string(filepath.Separator)) {
				return nil, fmt.Errorf("module %q icon path escapes addon directory", manifest.Name)
			}
		}
		manifests = append(manifests, manifest)
	}

	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Name < manifests[j].Name
	})
	return manifests, nil
}

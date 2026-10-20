package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/src"
	"github.com/evanw/esbuild/pkg/api"
)

var (
	watch           bool
	validModuleName = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
)

const (
	addonsDirectory       = "addons"
	staticDirectory       = "static"
	globalStaticDirectory = "web/web-app/resources/static"
)

type moduleManifest struct {
	Name string `json:"name"`
}

type assetManifest struct {
	Version int      `json:"version"`
	Scripts []string `json:"scripts"`
	Styles  []string `json:"styles"`
}

func main() {
	flagSet := newFlagSet()
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("failure", "error", err)
		os.Exit(1)
	}
}

func newFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("ts-compiler", flag.ContinueOnError)
	flagSet.BoolVar(&watch, "watch", false, "Enable watcher mode")
	flagSet.SetOutput(os.Stderr)
	return flagSet
}

func run(ctx context.Context) error {
	sctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	return build(sctx)
}

func build(ctx context.Context) error {
	legacyEntries, err := discoverLegacyEntries(src.TsDirectoryPath)
	if err != nil {
		return err
	}
	modules, err := discoverModuleAssets(addonsDirectory)
	if err != nil {
		return err
	}
	entries := append([]api.EntryPoint(nil), legacyEntries...)
	for _, module := range modules {
		entries = append(entries, module.entries...)
	}
	if err := cleanModuleAssetOutputs(addonsDirectory); err != nil {
		return err
	}
	if err := ensureModuleAssetMarkers(modules); err != nil {
		return err
	}
	if err := syncModulePublicAssets(modules); err != nil {
		return err
	}
	if len(entries) == 0 && !watch {
		slog.Info("no TypeScript or CSS entry points found")
		return nil
	}

	options := api.BuildOptions{
		EntryPointsAdvanced: legacyEntries,
		Bundle:              true,
		Format:              api.FormatESModule,
		LogLevel:            api.LogLevelInfo,
		MinifyIdentifiers:   true,
		MinifySyntax:        true,
		MinifyWhitespace:    true,
		Sourcemap:           api.SourceMapLinked,
		Target:              api.ESNext,
		Write:               true,
	}

	if watch {
		return watchBuild(ctx, options, modules)
	}
	if len(entries) == 0 {
		slog.Info("no TypeScript or CSS entry points found")
		return nil
	}

	slog.Info("building module assets", "entryPoints", len(entries))
	if len(legacyEntries) > 0 {
		legacyOptions := options
		legacyOptions.Outdir = globalStaticDirectory
		result := api.Build(legacyOptions)
		if err := buildError(result.Errors); err != nil {
			return err
		}
	}
	for _, module := range modules {
		if len(module.entries) == 0 {
			continue
		}
		moduleOptions := options
		moduleOptions.Outdir = module.distDir
		moduleOptions.EntryPointsAdvanced = module.entries
		result := api.Build(moduleOptions)
		if err := buildError(result.Errors); err != nil {
			return fmt.Errorf("build assets for module %q: %w", module.name, err)
		}
	}
	return nil
}

func watchBuild(ctx context.Context, options api.BuildOptions, modules []moduleAssetFiles) error {
	slog.Info("watching module assets", "entryPoints", len(options.EntryPointsAdvanced))
	contexts := make([]api.BuildContext, 0, len(modules)+1)
	if len(options.EntryPointsAdvanced) > 0 {
		legacyOptions := options
		legacyOptions.Outdir = globalStaticDirectory
		buildCtx, err := api.Context(legacyOptions)
		if err != nil {
			return err
		}
		if err := buildCtx.Watch(api.WatchOptions{}); err != nil {
			buildCtx.Dispose()
			return err
		}
		contexts = append(contexts, buildCtx)
	}
	for _, module := range modules {
		if len(module.entries) == 0 {
			continue
		}
		moduleOptions := options
		moduleOptions.Outdir = module.distDir
		moduleOptions.EntryPointsAdvanced = module.entries
		buildCtx, err := api.Context(moduleOptions)
		if err != nil {
			for _, ctx := range contexts {
				ctx.Dispose()
			}
			return err
		}
		if err := buildCtx.Watch(api.WatchOptions{}); err != nil {
			buildCtx.Dispose()
			for _, ctx := range contexts {
				ctx.Dispose()
			}
			return err
		}
		contexts = append(contexts, buildCtx)
	}
	defer func() {
		for _, ctx := range contexts {
			ctx.Dispose()
		}
	}()
	publicSignature, err := modulePublicSignature(modules)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			nextSignature, err := modulePublicSignature(modules)
			if err != nil {
				slog.Error("inspect module public assets", "error", err)
				continue
			}
			if nextSignature == publicSignature {
				continue
			}
			if err := syncModulePublicAssets(modules); err != nil {
				slog.Error("sync module public assets", "error", err)
				continue
			}
			publicSignature = nextSignature
			slog.Info("synced changed module public assets")
		}
	}
}

func buildError(messages []api.Message) error {
	if len(messages) == 0 {
		return nil
	}
	errs := make([]error, len(messages))
	for i, message := range messages {
		errs[i] = errors.New(message.Text)
	}
	return errors.Join(errs...)
}

func discoverEntries(addonsDir, legacyTsDir string) ([]api.EntryPoint, []moduleAssetFiles, error) {
	entries, err := discoverLegacyEntries(legacyTsDir)
	if err != nil {
		return nil, nil, err
	}

	modules, err := discoverModuleAssets(addonsDir)
	if err != nil {
		return nil, nil, err
	}
	for _, module := range modules {
		entries = append(entries, module.entries...)
	}
	return entries, modules, nil
}

func discoverLegacyEntries(tsDir string) ([]api.EntryPoint, error) {
	var entries []api.EntryPoint
	err := filepath.WalkDir(tsDir, func(filePath string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) && filePath == tsDir {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() || !isTypeScript(filePath) {
			return nil
		}

		relative, err := filepath.Rel(tsDir, filePath)
		if err != nil {
			return err
		}
		output := strings.TrimSuffix(filepath.ToSlash(relative), filepath.Ext(relative))
		entries = append(entries, api.EntryPoint{
			InputPath:  filePath,
			OutputPath: path.Join("js", output),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover legacy TypeScript files in %s: %w", tsDir, err)
	}
	return entries, nil
}

type moduleAssetFiles struct {
	name      string
	publicDir string
	distDir   string
	entries   []api.EntryPoint
}

func discoverModuleAssets(addonsDir string) ([]moduleAssetFiles, error) {
	addonEntries, err := os.ReadDir(addonsDir)
	if err != nil {
		return nil, fmt.Errorf("read addons directory %s: %w", addonsDir, err)
	}

	var modules []moduleAssetFiles
	for _, addon := range addonEntries {
		if !addon.IsDir() {
			continue
		}
		addonPath := filepath.Join(addonsDir, addon.Name())
		staticDir := filepath.Join(addonPath, staticDirectory)
		info, err := os.Stat(staticDir)
		if os.IsNotExist(err) {
			// Existing modules are migrated to the static asset contract gradually.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("inspect %s: %w", staticDir, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("%s must be a directory", staticDir)
		}

		module, err := loadModuleAssets(addonPath, addon.Name(), staticDir)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func loadModuleAssets(addonPath, directoryName, staticDir string) (moduleAssetFiles, error) {
	manifestPath := filepath.Join(addonPath, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return moduleAssetFiles{}, fmt.Errorf("read module manifest %s: %w", manifestPath, err)
	}
	var module moduleManifest
	if err := json.Unmarshal(manifestData, &module); err != nil {
		return moduleAssetFiles{}, fmt.Errorf("decode module manifest %s: %w", manifestPath, err)
	}
	if !validModuleName.MatchString(module.Name) {
		return moduleAssetFiles{}, fmt.Errorf("module manifest %s has invalid name %q", manifestPath, module.Name)
	}
	if module.Name != directoryName {
		return moduleAssetFiles{}, fmt.Errorf("module directory %q does not match manifest name %q", directoryName, module.Name)
	}

	assetsPath := filepath.Join(staticDir, "assets.json")
	assetsData, err := os.ReadFile(assetsPath)
	if err != nil {
		return moduleAssetFiles{}, fmt.Errorf("read assets manifest %s: %w", assetsPath, err)
	}
	assets, err := decodeAssetManifest(assetsData, assetsPath)
	if err != nil {
		return moduleAssetFiles{}, err
	}

	result := moduleAssetFiles{
		name:      module.Name,
		publicDir: filepath.Join(staticDir, "public"),
		distDir:   filepath.Join(staticDir, "dist"),
	}
	seen := make(map[string]struct{}, len(assets.Scripts)+len(assets.Styles))
	for _, script := range assets.Scripts {
		entry, key, err := moduleEntry(staticDir, script, "src/ts", "js", isTypeScript)
		if err != nil {
			return moduleAssetFiles{}, fmt.Errorf("module %q scripts: %w", module.Name, err)
		}
		if _, exists := seen[key]; exists {
			return moduleAssetFiles{}, fmt.Errorf("module %q declares asset entry %q more than once", module.Name, script)
		}
		seen[key] = struct{}{}
		result.entries = append(result.entries, entry)
	}
	for _, style := range assets.Styles {
		entry, key, err := moduleEntry(staticDir, style, "src/css", "css", func(p string) bool {
			return strings.EqualFold(filepath.Ext(p), ".css")
		})
		if err != nil {
			return moduleAssetFiles{}, fmt.Errorf("module %q styles: %w", module.Name, err)
		}
		if _, exists := seen[key]; exists {
			return moduleAssetFiles{}, fmt.Errorf("module %q declares asset entry %q more than once", module.Name, style)
		}
		seen[key] = struct{}{}
		result.entries = append(result.entries, entry)
	}
	return result, nil
}

func decodeAssetManifest(data []byte, filePath string) (assetManifest, error) {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var manifest assetManifest
	if err := decoder.Decode(&manifest); err != nil {
		return assetManifest{}, fmt.Errorf("decode assets manifest %s: %w", filePath, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return assetManifest{}, fmt.Errorf("decode assets manifest %s: multiple JSON values are not allowed", filePath)
		}
		return assetManifest{}, fmt.Errorf("decode assets manifest %s: %w", filePath, err)
	}
	if manifest.Version != 1 {
		return assetManifest{}, fmt.Errorf("assets manifest %s has unsupported version %d; expected 1", filePath, manifest.Version)
	}
	return manifest, nil
}

func moduleEntry(staticDir, relativePath, requiredRoot, outputRoot string, validExt func(string) bool) (api.EntryPoint, string, error) {
	if relativePath == "" || strings.Contains(relativePath, `\`) {
		return api.EntryPoint{}, "", fmt.Errorf("invalid path %q; use a non-empty slash-separated path relative to static/", relativePath)
	}
	clean := path.Clean(relativePath)
	if path.IsAbs(relativePath) || clean != relativePath || clean == "." || strings.HasPrefix(clean, "../") {
		return api.EntryPoint{}, "", fmt.Errorf("path %q must be clean and remain inside static/", relativePath)
	}
	if !strings.HasPrefix(clean, requiredRoot+"/") {
		return api.EntryPoint{}, "", fmt.Errorf("path %q must be inside %s/", relativePath, requiredRoot)
	}
	if !validExt(clean) {
		return api.EntryPoint{}, "", fmt.Errorf("unsupported entry point extension for %q", relativePath)
	}

	staticAbs, err := filepath.Abs(staticDir)
	if err != nil {
		return api.EntryPoint{}, "", err
	}
	inputAbs, err := filepath.Abs(filepath.Join(staticAbs, filepath.FromSlash(clean)))
	if err != nil {
		return api.EntryPoint{}, "", err
	}
	staticReal, err := filepath.EvalSymlinks(staticAbs)
	if err != nil {
		return api.EntryPoint{}, "", err
	}
	inputReal, err := filepath.EvalSymlinks(inputAbs)
	if err != nil {
		return api.EntryPoint{}, "", fmt.Errorf("resolve entry point %q: %w", relativePath, err)
	}
	withinStatic, err := filepath.Rel(staticReal, inputReal)
	if err != nil || withinStatic == ".." || strings.HasPrefix(withinStatic, ".."+string(filepath.Separator)) {
		return api.EntryPoint{}, "", fmt.Errorf("path %q escapes module static directory", relativePath)
	}
	info, err := os.Stat(inputAbs)
	if err != nil {
		return api.EntryPoint{}, "", fmt.Errorf("entry point %q: %w", relativePath, err)
	}
	if !info.Mode().IsRegular() {
		return api.EntryPoint{}, "", fmt.Errorf("entry point %q is not a regular file", relativePath)
	}

	insideRoot := strings.TrimPrefix(clean, requiredRoot+"/")
	output := path.Join(outputRoot, strings.TrimSuffix(insideRoot, path.Ext(insideRoot)))
	return api.EntryPoint{InputPath: inputAbs, OutputPath: output}, clean, nil
}

func isTypeScript(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".ts" || ext == ".tsx"
}

func copyModulePublicAssets(modules []moduleAssetFiles) error {
	for _, module := range modules {
		if _, err := os.Stat(module.publicDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("inspect public assets for module %q: %w", module.name, err)
		}

		err := filepath.WalkDir(module.publicDir, func(source string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Name() == ".gitkeep" {
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symbolic links are not supported in module %q public assets: %s", module.name, source)
			}
			if entry.IsDir() {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("public asset is not a regular file: %s", source)
			}

			relative, err := filepath.Rel(module.publicDir, source)
			if err != nil {
				return err
			}
			destination := filepath.Join(module.distDir, "public", relative)
			if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
				return fmt.Errorf("create public asset output directory: %w", err)
			}
			data, err := os.ReadFile(source)
			if err != nil {
				return fmt.Errorf("read public asset %s: %w", source, err)
			}
			if err := os.WriteFile(destination, data, info.Mode().Perm()); err != nil {
				return fmt.Errorf("write public asset %s: %w", destination, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("copy public assets for module %q: %w", module.name, err)
		}
	}
	return nil
}

func ensureModuleAssetMarkers(modules []moduleAssetFiles) error {
	for _, module := range modules {
		if err := os.MkdirAll(module.distDir, 0o755); err != nil {
			return fmt.Errorf("create generated asset directory for module %q: %w", module.name, err)
		}
		marker := filepath.Join(module.distDir, "module-assets.json")
		if err := os.WriteFile(marker, []byte("{\"generatedBy\":\"goep\"}\n"), 0o644); err != nil {
			return fmt.Errorf("write generated asset marker for module %q: %w", module.name, err)
		}
	}
	return nil
}

// cleanModuleAssetOutputs removes only generated dist directories below addons.
// Module sources and manifests under static/ are never touched here.
func cleanModuleAssetOutputs(addonsDir string) error {
	addonsAbs, err := filepath.Abs(addonsDir)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(addonsAbs)
	if err != nil {
		return fmt.Errorf("read addons directory %s: %w", addonsAbs, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || !validModuleName.MatchString(entry.Name()) {
			continue
		}
		staticAbs, err := filepath.Abs(filepath.Join(addonsAbs, entry.Name(), staticDirectory))
		if err != nil {
			return err
		}
		staticInfo, err := os.Lstat(staticAbs)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("inspect module static directory %s: %w", staticAbs, err)
		}
		if staticInfo.Mode()&os.ModeSymlink != 0 || !staticInfo.IsDir() {
			return fmt.Errorf("refusing to clean generated assets below non-directory or symlink %s", staticAbs)
		}
		distAbs, err := filepath.Abs(filepath.Join(staticAbs, "dist"))
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(staticAbs, distAbs)
		if err != nil || relative != "dist" {
			return fmt.Errorf("refusing to clean unexpected module asset output %q", distAbs)
		}
		if err := os.RemoveAll(distAbs); err != nil {
			return fmt.Errorf("clean generated module assets %s: %w", distAbs, err)
		}
	}
	return nil
}

func syncModulePublicAssets(modules []moduleAssetFiles) error {
	for _, module := range modules {
		publicOutput := filepath.Join(module.distDir, "public")
		if err := os.RemoveAll(publicOutput); err != nil {
			return fmt.Errorf("clean public asset output for module %q: %w", module.name, err)
		}
	}
	return copyModulePublicAssets(modules)
}

func modulePublicSignature(modules []moduleAssetFiles) (string, error) {
	var files []string
	for _, module := range modules {
		err := filepath.WalkDir(module.publicDir, func(filePath string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				if os.IsNotExist(walkErr) && filePath == module.publicDir {
					return nil
				}
				return walkErr
			}
			if entry.IsDir() || entry.Name() == ".gitkeep" {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			files = append(files, fmt.Sprintf("%s|%d|%d", filePath, info.Size(), info.ModTime().UnixNano()))
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("inspect public assets for module %q: %w", module.name, err)
		}
	}
	sort.Strings(files)
	return strings.Join(files, "\n"), nil
}

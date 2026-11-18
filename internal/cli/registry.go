package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type manifest struct {
	Type              string   `json:"type,omitempty"`
	Name              string   `json:"name"`
	Depends           []string `json:"depends"`
	DatabaseProviders []string `json:"database_providers,omitempty"`
}

type moduleInfo struct {
	name, importPath, assetsImportPath, migrationsImportPath string
	kind                                                     string
	depends                                                  []string
	hasStaticAssets, hasPostgresMigrations                   bool
}

func GenerateModuleRegistry(args []string) error {
	flags := flag.NewFlagSet("generate-module-registry", flag.ContinueOnError)
	addonsDir := flags.String("addons", "addons", "directory containing module manifests")
	output := flags.String("out", "internal/modules/backend/factories_gen.go", "generated Go source file")
	assetsOutput := flags.String("assets-out", "internal/modules/assets/static_assets_gen.go", "generated module static asset registry")
	migrationsOutput := flags.String("migrations-out", "internal/modules/migrations/addons_gen.go", "generated addon PostgreSQL migration registry")
	modulePath := flags.String("module", "github.com/GoEnterpricePlatform/goEP-core", "Go module import path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	return generate(*addonsDir, *output, *assetsOutput, *migrationsOutput, *modulePath)
}

func generate(addonsDir, output, assetsOutput, migrationsOutput, modulePath string) error {
	entries, err := os.ReadDir(addonsDir)
	if err != nil {
		return fmt.Errorf("read addons directory: %w", err)
	}
	var modules []moduleInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		addonDir := filepath.Join(addonsDir, entry.Name())
		data, err := os.ReadFile(filepath.Join(addonDir, "manifest.json"))
		if err != nil {
			return fmt.Errorf("read %s manifest: %w", entry.Name(), err)
		}
		var m manifest
		if err := json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("decode %s manifest: %w", entry.Name(), err)
		}
		if m.Type != "" && m.Type != "backend" && m.Type != "web" {
			return fmt.Errorf("%s manifest has unsupported module type %q", entry.Name(), m.Type)
		}
		if m.Name == "" {
			return fmt.Errorf("%s manifest has no module name", entry.Name())
		}
		configDir := filepath.Join(addonDir, "config")
		registeredMethod := "RegisterBackendModule"
		if m.Type == "web" {
			registeredMethod = "RegisterWebModule"
			registered, err := hasRegistration(configDir, registeredMethod)
			if err != nil {
				return fmt.Errorf("module %q: %w", m.Name, err)
			}
			if !registered {
				return fmt.Errorf("module %q: config package must define %s(*modules.Registry) error", m.Name, registeredMethod)
			}
		} else {
			if _, err = findConstructor(configDir); err != nil {
				return fmt.Errorf("module %q: %w", m.Name, err)
			}
			registered, err := hasRegistration(configDir, registeredMethod)
			if err != nil {
				return fmt.Errorf("module %q: %w", m.Name, err)
			}
			if !registered {
				return fmt.Errorf("module %q: config package must define %s(*modules.Registry) error", m.Name, registeredMethod)
			}
		}
		addonImportPath := strings.TrimSuffix(modulePath, "/") + "/addons/" + filepath.ToSlash(entry.Name())
		info := moduleInfo{
			kind:       m.Type,
			name:       m.Name,
			depends:    append([]string(nil), m.Depends...),
			importPath: addonImportPath + "/config",
		}
		staticDir := filepath.Join(addonDir, "static")
		if staticInfo, err := os.Stat(staticDir); err == nil && staticInfo.IsDir() {
			info.assetsImportPath = addonImportPath
			info.hasStaticAssets = true
			if err := writeModuleAssetsSource(addonDir); err != nil {
				return fmt.Errorf("generate static assets package for module %q: %w", m.Name, err)
			}
		} else if err == nil {
			return fmt.Errorf("module %q static path is not a directory", m.Name)
		} else if os.IsNotExist(err) {
			generatedSource := filepath.Join(addonDir, "static_assets_gen.go")
			if removeErr := os.Remove(generatedSource); removeErr != nil && !os.IsNotExist(removeErr) {
				return fmt.Errorf("remove obsolete static assets package for module %q: %w", m.Name, removeErr)
			}
		} else if err != nil {
			return fmt.Errorf("inspect static directory for module %q: %w", m.Name, err)
		}
		migrationDir := filepath.Join(addonDir, "database", "postgresql", "migrations")
		if m.Type != "web" && supportsPostgres(m.DatabaseProviders) {
			upFiles, err := filepath.Glob(filepath.Join(migrationDir, "*.up.sql"))
			if err != nil {
				return fmt.Errorf("inspect PostgreSQL migrations for module %q: %w", m.Name, err)
			}
			downFiles, err := filepath.Glob(filepath.Join(migrationDir, "*.down.sql"))
			if err != nil {
				return fmt.Errorf("inspect PostgreSQL down migrations for module %q: %w", m.Name, err)
			}
			if len(upFiles) == 0 && len(downFiles) > 0 {
				return fmt.Errorf("module %q has down migrations but no up migrations", m.Name)
			}
			if len(upFiles) > 0 {
				info.migrationsImportPath = addonImportPath + "/database/postgresql/migrations"
				info.hasPostgresMigrations = true
				if err := writeModuleMigrationsSource(migrationDir); err != nil {
					return fmt.Errorf("generate PostgreSQL migrations package for module %q: %w", m.Name, err)
				}
			}
		} else {
			generatedSource := filepath.Join(migrationDir, "migrations_gen.go")
			if removeErr := os.Remove(generatedSource); removeErr != nil && !os.IsNotExist(removeErr) {
				return fmt.Errorf("remove obsolete PostgreSQL migrations package for module %q: %w", m.Name, removeErr)
			}
		}
		modules = append(modules, info)
	}
	sort.Slice(modules, func(i, j int) bool { return modules[i].name < modules[j].name })

	var source bytes.Buffer
	fmt.Fprintln(&source, "// Code generated by goep generate-module-registry; DO NOT EDIT.")
	fmt.Fprintln(&source, "package backend")
	fmt.Fprintln(&source, "import (")
	fmt.Fprintf(&source, "\t\"%s/internal/modules\"\n", strings.TrimSuffix(modulePath, "/"))
	for i, m := range modules {
		fmt.Fprintf(&source, "\tmodule%d %q\n", i, m.importPath)
	}
	fmt.Fprintln(&source, ")")
	fmt.Fprintln(&source, "// NewCompiledRegistry registers backend and web module initializers.")
	fmt.Fprintln(&source, "func NewCompiledRegistry() (*modules.Registry, error) {")
	fmt.Fprintln(&source, "\tregistry := modules.NewRegistry()")
	for i, m := range modules {
		if m.kind == "web" {
			fmt.Fprintf(&source, "\tif err := module%d.RegisterWebModule(registry); err != nil { return nil, err }\n", i)
		} else {
			fmt.Fprintf(&source, "\tif err := module%d.RegisterBackendModule(registry); err != nil { return nil, err }\n", i)
		}
	}
	fmt.Fprintln(&source, "\treturn registry, nil")
	fmt.Fprintln(&source, "}")
	formatted, err := format.Source(source.Bytes())
	if err != nil {
		return fmt.Errorf("format generated registry: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := writeFileIfChanged(output, formatted, 0o644); err != nil {
		return fmt.Errorf("write generated registry: %w", err)
	}
	if err := writeStaticAssetsRegistry(assetsOutput, modules); err != nil {
		return fmt.Errorf("write generated static assets registry: %w", err)
	}
	if err := writeMigrationsRegistry(migrationsOutput, modules); err != nil {
		return fmt.Errorf("write generated PostgreSQL migrations registry: %w", err)
	}
	return nil
}

func supportsPostgres(providers []string) bool {
	if len(providers) == 0 {
		return true
	}
	for _, provider := range providers {
		if strings.TrimSpace(provider) == "postgresql" {
			return true
		}
	}
	return false
}

func writeModuleMigrationsSource(migrationDir string) error {
	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		return err
	}
	const source = `// Code generated by goep generate-module-registry; DO NOT EDIT.
package migrations

import "embed"

// Files embeds this addon's PostgreSQL migrations.
//
//go:embed *.sql
var Files embed.FS
`
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return err
	}
	return writeFileIfChanged(filepath.Join(migrationDir, "migrations_gen.go"), formatted, 0o644)
}

func writeModuleAssetsSource(addonDir string) error {
	staticDir := filepath.Join(addonDir, "static")
	distDir := filepath.Join(staticDir, "dist")
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return err
	}
	marker := filepath.Join(distDir, "module-assets.json")
	if err := writeFileIfChanged(marker, []byte("{\"generatedBy\":\"goep\"}\n"), 0o644); err != nil {
		return err
	}

	const source = `// Code generated by goep generate-module-registry; DO NOT EDIT.
package module

import "embed"

// StaticAssets embeds this addon's generated Web App files.
//
//go:embed static/dist
var StaticAssets embed.FS
`
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return err
	}
	return writeFileIfChanged(filepath.Join(addonDir, "static_assets_gen.go"), formatted, 0o644)
}

func writeFileIfChanged(path string, contents []byte, mode os.FileMode) error {
	current, err := os.ReadFile(path)
	if err == nil {
		// Generated files may already use CRLF on Windows. Compare normalized
		// line endings so an otherwise identical file is not rewritten on each run.
		normalizedCurrent := bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n"))
		if bytes.Equal(normalizedCurrent, contents) {
			return nil
		}
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, contents, mode)
}

func writeStaticAssetsRegistry(output string, modules []moduleInfo) error {
	var source bytes.Buffer
	fmt.Fprintln(&source, "// Code generated by goep generate-module-registry; DO NOT EDIT.")
	fmt.Fprintln(&source, "package assets")
	fmt.Fprintln(&source, "import (")
	fmt.Fprintln(&source, "\t\"io/fs\"")
	if hasModuleStaticAssets(modules) {
		fmt.Fprintln(&source)
	}
	for i, module := range modules {
		if module.hasStaticAssets {
			fmt.Fprintf(&source, "\tassetModule%d %q\n", i, module.assetsImportPath)
		}
	}
	fmt.Fprintln(&source, ")")
	fmt.Fprintln(&source, "// StaticDirectories maps addon names to their embedded static filesystems.")
	fmt.Fprintln(&source, "var StaticDirectories = map[string]fs.FS{")
	for i, module := range modules {
		if module.hasStaticAssets {
			fmt.Fprintf(&source, "\t%q: assetModule%d.StaticAssets,\n", module.name, i)
		}
	}
	fmt.Fprintln(&source, "}")
	formatted, err := format.Source(source.Bytes())
	if err != nil {
		return fmt.Errorf("format generated static assets registry: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create static registry output directory: %w", err)
	}
	if err := writeFileIfChanged(output, formatted, 0o644); err != nil {
		return fmt.Errorf("write static registry: %w", err)
	}
	return nil
}

func writeMigrationsRegistry(output string, modules []moduleInfo) error {
	var source bytes.Buffer
	fmt.Fprintln(&source, "// Code generated by goep generate-module-registry; DO NOT EDIT.")
	fmt.Fprintln(&source, "package migrations")
	fmt.Fprintln(&source, "import (")
	for i, module := range modules {
		if module.hasPostgresMigrations {
			fmt.Fprintf(&source, "\taddonMigration%d %q\n", i, module.migrationsImportPath)
		}
	}
	fmt.Fprintln(&source, ")")
	fmt.Fprintln(&source, "// addonMigrations contains PostgreSQL migration files discovered from addon manifests.")
	fmt.Fprintln(&source, "var addonMigrations = []moduleMigrations{")
	for i, module := range modules {
		if !module.hasPostgresMigrations {
			continue
		}
		depends := make([]string, len(module.depends))
		for j, dependency := range module.depends {
			depends[j] = fmt.Sprintf("%q", dependency)
		}
		fmt.Fprintf(&source, "\t{name: %q, fs: addonMigration%d.Files, depends: []string{%s}},\n", module.name, i, strings.Join(depends, ", "))
	}
	fmt.Fprintln(&source, "}")
	formatted, err := format.Source(source.Bytes())
	if err != nil {
		return fmt.Errorf("format generated PostgreSQL migrations registry: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create PostgreSQL migrations registry directory: %w", err)
	}
	if err := writeFileIfChanged(output, formatted, 0o644); err != nil {
		return fmt.Errorf("write PostgreSQL migrations registry: %w", err)
	}
	return nil
}

func hasModuleStaticAssets(modules []moduleInfo) bool {
	for _, module := range modules {
		if module.hasStaticAssets {
			return true
		}
	}
	return false
}

func hasRegistration(configDir, functionName string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(configDir, "*.go"))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return false, fmt.Errorf("parse %s: %w", file, err)
		}
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && fn.Name.Name == functionName {
				return true, nil
			}
		}
	}
	return false, nil
}

func findConstructor(configDir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(configDir, "*.go"))
	if err != nil {
		return "", err
	}
	fset := token.NewFileSet()
	var matches []string
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", file, err)
		}
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !ast.IsExported(fn.Name.Name) || fn.Type.Params == nil || fn.Type.Results == nil {
				continue
			}
			if len(fn.Type.Params.List) != 1 || len(fn.Type.Results.List) != 2 {
				continue
			}
			if typeName(fn.Type.Params.List[0].Type) != "ModuleConfig" || typeName(fn.Type.Results.List[0].Type) != "*Module" || typeName(fn.Type.Results.List[1].Type) != "error" {
				continue
			}
			matches = append(matches, fn.Name.Name)
		}
	}
	if len(matches) != 1 {
		return "", fmt.Errorf("expected exactly one exported constructor func(Config ModuleConfig) (*Module, error), found %d", len(matches))
	}
	return matches[0], nil
}

func typeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeName(t.X)
	default:
		return ""
	}
}

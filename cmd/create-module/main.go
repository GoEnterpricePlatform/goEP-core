package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const modulePath = "github.com/GoEnterpricePlatform/goEP-core"

var validModuleName = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: go run ./cmd/create-module [--summary \"description\"] <module-name>\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Example: go run ./cmd/create-module --summary \"Sales reports\" sales-reports")
		flag.PrintDefaults()
	}
	summary := flag.String("summary", "", "short description written to the module manifest")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	name := flag.Arg(0)
	if !validModuleName.MatchString(name) {
		fatalf("invalid module name %q: use lowercase letters and numbers separated by single hyphens", name)
	}

	root, err := findRepositoryRoot()
	if err != nil {
		fatalf("find repository root: %v", err)
	}
	if *summary == "" {
		*summary = strings.ReplaceAll(name, "-", " ") + " module"
	}

	moduleDir, err := createModule(root, name, *summary)
	if err != nil {
		fatalf("create module: %v", err)
	}
	fmt.Printf("Created %s\n", filepath.ToSlash(filepath.Join("addons", name)))
	fmt.Println("Generating the compiled module registry...")

	generate := exec.Command("go", "generate", "./internal/modules/backend")
	generate.Dir = root
	generate.Stdout = os.Stdout
	generate.Stderr = os.Stderr
	if err := generate.Run(); err != nil {
		fatalf("module files were created at %s, but registry generation failed: %v", moduleDir, err)
	}
	fmt.Printf("Module %q is ready and included in the compiled registry.\n", name)
}

func createModule(root, name, summary string) (string, error) {
	addonsDir := filepath.Join(root, "addons")
	moduleDir := filepath.Join(addonsDir, name)
	if _, err := os.Stat(moduleDir); err == nil {
		return "", fmt.Errorf("%s already exists", moduleDir)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("check %s: %w", moduleDir, err)
	}

	tempDir, err := os.MkdirTemp(addonsDir, ".create-module-")
	if err != nil {
		return "", fmt.Errorf("create temporary module directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	identifier := goIdentifier(name)
	files := map[string]string{
		"manifest.json": manifestTemplate(name, summary),
		filepath.Join("config", "module_config.go"):  moduleConfigTemplate(identifier),
		filepath.Join("config", "backend_module.go"): backendModuleTemplate(name),
		filepath.Join("domain", "domain.go"):         domainTemplate(name),
	}
	for relativePath, contents := range files {
		path := filepath.Join(tempDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", fmt.Errorf("create directory for %s: %w", relativePath, err)
		}
		if strings.HasSuffix(path, ".go") {
			formatted, err := format.Source([]byte(contents))
			if err != nil {
				return "", fmt.Errorf("format %s: %w", relativePath, err)
			}
			contents = string(formatted)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			return "", fmt.Errorf("write %s: %w", relativePath, err)
		}
	}

	if err := os.Rename(tempDir, moduleDir); err != nil {
		return "", fmt.Errorf("install module directory: %w", err)
	}
	return moduleDir, nil
}

func findRepositoryRoot() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(current, "addons")); err == nil {
				return current, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("could not find go.mod and addons/; run this command inside the goEP-core repository")
		}
		current = parent
	}
}

func goIdentifier(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func manifestTemplate(name, summary string) string {
	return fmt.Sprintf(`{
  "name": %q,
  "version": "1.0.0",
  "summary": %q,
  "depends": []
}
`, name, summary)
}

func moduleConfigTemplate(identifier string) string {
	return fmt.Sprintf(`package config

import (
	"net/http"

	appConfig %q
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	DB         *mongo.Database
}

type Module struct{}

func New%sModule(_ ModuleConfig) (*Module, error) {
	return &Module{}, nil
}
`, modulePath+"/internal/config", identifier)
}

func backendModuleTemplate(name string) string {
	identifier := goIdentifier(name)
	return fmt.Sprintf(`package config

import "github.com/GoEnterpricePlatform/goEP-core/internal/modules"

func RegisterBackendModule(registry *modules.Registry) error {
	if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
		Name: %q,
		New:  New%sModule,
	}); err != nil {
		return err
	}

	return modules.RegisterBackendInitializer(registry, %q, func(runtime modules.BackendRuntime) (any, error) {
		return New%sModule(ModuleConfig{
			AppEnvs:    runtime.AppEnvs,
			AppClients: runtime.AppClients,
			APIv1:      runtime.APIv1,
			DB:         runtime.DB,
		})
	})
}
`, name, identifier, name, identifier)
}

func domainTemplate(name string) string {
	return fmt.Sprintf(`// Package domain contains the business types and rules for the %s module.
package domain
`, strings.ReplaceAll(name, "-", " "))
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "create-module: "+format+"\n", args...)
	os.Exit(1)
}

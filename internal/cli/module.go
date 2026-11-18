package cli

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const modulePath = "github.com/GoEnterpricePlatform/goEP-core"

var validModuleName = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func CreateModule(args []string) error {
	flags := flag.NewFlagSet("create-module", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: go run ./cmd/goep create-module [--summary \"description\"] [--database postgresql|mongo|both|none] [--ui simple|dashboard] [--ui-layout sidebar|replace] <module-name>")
		fmt.Fprintln(flags.Output(), "Example: go run ./cmd/goep create-module --summary \"Sales reports\" sales-reports")
		flags.PrintDefaults()
	}
	summary := flags.String("summary", "", "short description written to the module manifest")
	database := flags.String("database", "", "proveedor de base de datos: postgresql, mongo, both o none (predeterminado: preguntar)")
	ui := flags.String("ui", "", "tipo de interfaz: simple o dashboard (predeterminado: preguntar)")
	uiLayout := flags.String("ui-layout", "", "presentación: sidebar o replace (predeterminado: preguntar)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return fmt.Errorf("expected one module name")
	}
	name := flags.Arg(0)
	if !validModuleName.MatchString(name) {
		return fmt.Errorf("invalid module name %q: use lowercase letters and numbers separated by single hyphens", name)
	}
	input := bufio.NewReader(os.Stdin)
	provider, err := chooseDatabaseProvider(*database, input, os.Stdout)
	if err != nil {
		return err
	}
	uiType, err := chooseUIType(*ui, input, os.Stdout)
	if err != nil {
		return err
	}
	uiPresentation, err := chooseUILayout(*uiLayout, uiType, input, os.Stdout)
	if err != nil {
		return err
	}

	root, err := findRepositoryRoot()
	if err != nil {
		return fmt.Errorf("find repository root: %w", err)
	}
	if *summary == "" {
		*summary = strings.ReplaceAll(name, "-", " ") + " module"
	}

	moduleDir, err := createModule(root, name, *summary, provider, uiType, uiPresentation)
	if err != nil {
		return fmt.Errorf("create module: %w", err)
	}
	fmt.Printf("Created %s\n", filepath.ToSlash(filepath.Join("addons", name)))
	fmt.Println("Generating the module UI templates...")
	templ := exec.Command("go", "tool", "templ", "generate", "./addons/"+name+"/ui")
	templ.Dir = root
	templ.Stdout = os.Stdout
	templ.Stderr = os.Stderr
	if err := templ.Run(); err != nil {
		return fmt.Errorf("module files were created at %s, but UI template generation failed: %w", moduleDir, err)
	}
	fmt.Println("Generating the backend and static asset registries...")

	generate := exec.Command("go", "generate", "./internal/modules/backend")
	generate.Dir = root
	generate.Stdout = os.Stdout
	generate.Stderr = os.Stderr
	if err := generate.Run(); err != nil {
		return fmt.Errorf("module files were created at %s, but registry generation failed: %w", moduleDir, err)
	}
	fmt.Printf("Module %q is ready and included in the backend and static asset registries.\n", name)
	return nil
}

func chooseUIType(value string, input io.Reader, output *os.File) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value != "" {
		switch value {
		case "simple", "dashboard":
			return value, nil
		case "sidebar":
			return "dashboard", nil
		default:
			return "", fmt.Errorf("unsupported UI type %q: choose simple or dashboard", value)
		}
	}

	fmt.Fprintln(output, "¿Qué tipo de interfaz tendrá el módulo?")
	fmt.Fprintln(output, "  1) Simple: una página")
	fmt.Fprintln(output, "  2) Dashboard (layout de Identity)")
	fmt.Fprint(output, "Selecciona [1]: ")
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", fmt.Errorf("leer tipo de interfaz: %w", err)
	}
	switch strings.TrimSpace(line) {
	case "", "1":
		return "simple", nil
	case "2":
		return "dashboard", nil
	default:
		return "", fmt.Errorf("selección de interfaz no válida %q: elige 1 o 2", strings.TrimSpace(line))
	}
}

func chooseUILayout(value, uiType string, input io.Reader, output *os.File) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value != "" {
		switch value {
		case "sidebar", "replace":
			return value, nil
		default:
			return "", fmt.Errorf("unsupported UI layout %q: choose sidebar or replace", value)
		}
	}

	defaultLayout := "replace"
	if uiType == "dashboard" {
		defaultLayout = "sidebar"
	}
	fmt.Fprintln(output, "¿Cómo se mostrará la pantalla del módulo?")
	fmt.Fprintln(output, "  1) Dentro del dashboard con sidebar de Identity")
	fmt.Fprintln(output, "  2) Reemplazar la pantalla del módulo")
	fmt.Fprintf(output, "Selecciona [%s]: ", map[string]string{"sidebar": "1", "replace": "2"}[defaultLayout])
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", fmt.Errorf("leer presentación de la interfaz: %w", err)
	}
	switch strings.TrimSpace(line) {
	case "":
		return defaultLayout, nil
	case "1", "sidebar":
		return "sidebar", nil
	case "2", "replace":
		return "replace", nil
	default:
		return "", fmt.Errorf("selección de presentación no válida %q: elige 1 o 2", strings.TrimSpace(line))
	}
}

func chooseDatabaseProvider(value string, input io.Reader, output *os.File) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value != "" {
		switch value {
		case "postgresql", "mongo", "both", "none":
			return value, nil
		default:
			return "", fmt.Errorf("unsupported database provider %q: choose postgresql, mongo, both, or none", value)
		}
	}

	fmt.Fprintln(output, "¿Qué base de datos usará el módulo?")
	fmt.Fprintln(output, "  1) PostgreSQL (default)")
	fmt.Fprintln(output, "  2) MongoDB")
	fmt.Fprintln(output, "  3) Both")
	fmt.Fprintln(output, "  4) none")
	fmt.Fprint(output, "Select [1]: ")
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", fmt.Errorf("read database provider: %w", err)
	}
	switch strings.TrimSpace(line) {
	case "", "1":
		return "postgresql", nil
	case "2":
		return "mongo", nil
	case "3":
		return "both", nil
	case "4":
		return "none", nil
	default:
		return "", fmt.Errorf("invalid database selection %q: choose 1, 2, 3, or 4", strings.TrimSpace(line))
	}
}

func createModule(root, name, summary, provider, uiType, uiLayout string) (string, error) {
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

	files := map[string]string{
		"manifest.json": manifestTemplate(name, summary, provider, uiType, uiLayout),
		filepath.Join("config", "module_config.go"):       moduleConfigTemplateWithUI(name),
		filepath.Join("config", "backend_module.go"):      backendModuleTemplate(name),
		filepath.Join("domain", "domain.go"):              domainTemplate(name),
		filepath.Join("static", "assets.json"):            assetsManifestTemplate(),
		filepath.Join("static", "src", "ts", ".gitkeep"):  "",
		filepath.Join("static", "src", "css", ".gitkeep"): "",
		filepath.Join("static", "public", ".gitkeep"):     "",
		filepath.Join("ui", "handler", "handler.go"):      uiHandlerTemplate(name, uiLayout),
		filepath.Join("ui", "pages", "page.templ"):        uiPageTemplate(name, uiType),
	}
	if provider == "postgresql" || provider == "both" {
		files[filepath.Join("database", "postgresql", "migrations", "000001_initial.up.sql")] = initialUpMigrationTemplate(name)
		files[filepath.Join("database", "postgresql", "migrations", "000001_initial.down.sql")] = initialDownMigrationTemplate(name)
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

func manifestTemplate(name, summary, provider, uiType, uiLayout string) string {
	providers := ""
	switch provider {
	case "postgresql":
		providers = `,
  "database_providers": ["postgresql"]`
	case "mongo":
		providers = `,
  "database_providers": ["mongo"]`
	case "both":
		providers = `,
  "database_providers": ["postgresql", "mongo"]`
	}
	return fmt.Sprintf(`{
  "name": %q,
  "version": "1.0.0",
  "summary": %q,
  "ui": {"style": %q, "layout": %q},
  "depends": []%s
}
`, name, summary, uiType, uiLayout, providers)
}

func moduleConfigTemplateWithUI(name string) string {
	identifier := goIdentifier(name)
	return fmt.Sprintf(`package config

import (
	"net/http"

	appConfig %q
	uiHandler %q
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	TemplateV1 *http.ServeMux
	DB         *mongo.Database
	PGConn     *pgx.Conn
}

type Module struct {
	UI *uiHandler.Handler
}

func New%sModule(cfg ModuleConfig) (*Module, error) {
	ui := uiHandler.NewHandler()
	ui.RegisterRoutes(cfg.TemplateV1)
	return &Module{UI: ui}, nil
}
`, modulePath+"/internal/config", modulePath+"/addons/"+name+"/ui/handler", identifier)
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
			TemplateV1: runtime.TemplateV1,
			DB:         runtime.DB,
			PGConn:     runtime.PGConn,
		})
	})
}
`, name, identifier, name, identifier)
}

func uiHandlerTemplate(name, uiLayout string) string {
	identifier := goIdentifier(name)
	return fmt.Sprintf(`package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/layouts"
	"github.com/GoEnterpricePlatform/goEP-core/addons/%s/ui/pages"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	mux.HandleFunc("GET /goep-admin/modules/%s", h.render)
	mux.HandleFunc("GET /%s", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/v1/%s/", http.StatusPermanentRedirect)
	})
	mux.HandleFunc("GET /%s/", h.render)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request) {
	var page = pages.%sPage()
	if %q == "sidebar" {
		page = layouts.DashboardLayout("overview", []string{"ROLE_SYSTEM_ADMIN"}, nil, []navigation.Item{{
			Title: "Overview", Href: "/v1/%s/", ActivePage: "overview",
		}}, pages.%sContent())
	}
	if err := page.Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
`, name, name, name, name, name, identifier, uiLayout, name, identifier)
}

func uiPageTemplate(name, uiType string) string {
	title := strings.Title(strings.ReplaceAll(name, "-", " "))
	identifier := goIdentifier(name)
	style := "simple"
	if uiType == "dashboard" {
		style = "dashboard"
	}
	return fmt.Sprintf(`package pages

import "github.com/GoEnterpricePlatform/goEP-core/addons/web-app/resources"

templ %sContent() {
	<section>
		<div class="mb-6">
			<h1 class="text-2xl font-bold tracking-tight text-black">%s</h1>
			<p class="mt-1 text-sm text-zinc-500">Your %s module is ready.</p>
		</div>
		if %q == "dashboard" {
			<div class="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-4">
				<div class="rounded-xl border border-zinc-200 bg-white p-5"><p class="text-sm text-zinc-500">Overview</p><p class="mt-2 text-2xl font-bold">Ready</p></div>
				<div class="rounded-xl border border-zinc-200 bg-white p-5"><p class="text-sm text-zinc-500">Module</p><p class="mt-2 text-2xl font-bold">%s</p></div>
			</div>
		}
	</section>
}

templ %sPage() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>%s</title>
			<link href={ resources.StaticPath("css/app.css") } rel="stylesheet" type="text/css"/>
		</head>
		<body class="min-h-screen bg-white text-zinc-900 antialiased">
			<main class="mx-auto max-w-7xl px-6 py-10 lg:px-8">
				@%sContent()
			</main>
		</body>
	</html>
}
`, identifier, title, style, style, title, identifier, title, identifier)
}

func domainTemplate(name string) string {
	return fmt.Sprintf(`// Package domain contains the business types and rules for the %s module.
package domain
`, strings.ReplaceAll(name, "-", " "))
}

func assetsManifestTemplate() string {
	return `{
  "version": 1,
  "scripts": [],
  "styles": []
}
`
}

func initialUpMigrationTemplate(name string) string {
	return fmt.Sprintf("-- Initial PostgreSQL migration for the %s module.\n-- Add this module's initial tables and indexes below.\n", strings.ReplaceAll(name, "-", " "))
}

func initialDownMigrationTemplate(name string) string {
	return fmt.Sprintf("-- Roll back the initial PostgreSQL migration for the %s module.\n-- Add DROP statements for objects created by 000001_initial.up.sql below.\n", strings.ReplaceAll(name, "-", " "))
}

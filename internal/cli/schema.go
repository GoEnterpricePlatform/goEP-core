// Command assemble-schema builds the PostgreSQL schema bundle from the core
// schema and the schemas owned by each addon.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type schemaFile struct {
	name string
	path string
}

var createTable = regexp.MustCompile(`(?im)^CREATE\s+TABLE\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
var createIndex = regexp.MustCompile(`(?im)^CREATE\s+(UNIQUE\s+)?INDEX\s+([a-zA-Z_][a-zA-Z0-9_]*)\b`)

func AssembleSchema() error {
	files := []schemaFile{
		{name: "identity", path: filepath.Join("addons", "identity", "database", "postgresql", "schema.sql")},
		{name: "posts", path: filepath.Join("addons", "posts", "database", "postgresql", "schema.sql")},
		{name: "catalog", path: filepath.Join("addons", "catalog", "database", "postgresql", "schema.sql")},
		{name: "billing", path: filepath.Join("addons", "billing", "database", "postgresql", "schema.sql")},
	}

	var schema bytes.Buffer
	schema.WriteString("-- GENERATED FILE. Do not edit; update the source schema files and run go run ./cmd/goep schema.\n\n")
	var tables []string
	for _, file := range files {
		contents, err := os.ReadFile(file.path)
		if err != nil {
			return fmt.Errorf("read %s schema: %w", file.name, err)
		}
		matches := createTable.FindAllSubmatch(contents, -1)
		if len(matches) == 0 {
			return fmt.Errorf("%s schema has no CREATE TABLE statements", file.name)
		}
		for _, match := range matches {
			tables = append(tables, string(match[1]))
		}
		contents = createTable.ReplaceAll(contents, []byte("CREATE TABLE IF NOT EXISTS $1 ("))
		contents = createIndex.ReplaceAll(contents, []byte("CREATE ${1}INDEX IF NOT EXISTS ${2}"))
		fmt.Fprintf(&schema, "-- Source: %s\n", filepath.ToSlash(file.path))
		schema.Write(bytes.TrimSpace(contents))
		schema.WriteString("\n\n")
	}

	if err := os.MkdirAll(filepath.Join("internal", "clients", "postgresql", "model"), 0o755); err != nil {
		return fmt.Errorf("create schema output directory: %w", err)
	}
	outputDir := filepath.Join("internal", "clients", "postgresql", "model")
	if err := os.WriteFile(filepath.Join(outputDir, "schema.sql"), schema.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write schema.sql: %w", err)
	}

	var reset bytes.Buffer
	reset.WriteString("-- GENERATED FILE. Do not edit; run go run ./cmd/goep schema.\n")
	reset.WriteString("-- Development only: this script drops all application tables.\n\n")
	for i := len(tables) - 1; i >= 0; i-- {
		fmt.Fprintf(&reset, "DROP TABLE IF EXISTS %s CASCADE;\n", tables[i])
	}
	reset.WriteString("\n")
	reset.Write(schema.Bytes())
	if err := os.WriteFile(filepath.Join(outputDir, "dev_reset.sql"), reset.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write dev_reset.sql: %w", err)
	}

	managedMigrationsDir := filepath.Join("internal", "modules", "backend", "managed", "migrations")
	managedMigrations, err := filepath.Glob(filepath.Join(managedMigrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("list managed lifecycle migrations: %w", err)
	}
	if len(managedMigrations) == 0 {
		return fmt.Errorf("no managed lifecycle up migrations found in %s", managedMigrationsDir)
	}
	sort.Strings(managedMigrations)
	var managedBundle bytes.Buffer
	managedBundle.WriteString("-- MANAGED FLAVOR ONLY. Generated from managed up migrations.\n\n")
	for _, migrationPath := range managedMigrations {
		contents, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("read managed lifecycle migration %s: %w", migrationPath, err)
		}
		fmt.Fprintf(&managedBundle, "-- Source: %s\n", filepath.ToSlash(migrationPath))
		managedBundle.Write(bytes.TrimSpace(contents))
		managedBundle.WriteString("\n\n")
	}
	managedBundleBytes := append(bytes.TrimRight(managedBundle.Bytes(), "\n"), '\n')
	if err := os.WriteFile(filepath.Join(outputDir, "managed_schema.sql"), managedBundleBytes, 0o644); err != nil {
		return fmt.Errorf("write managed_schema.sql: %w", err)
	}
	return nil
}

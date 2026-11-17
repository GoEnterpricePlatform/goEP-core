// Package migrations applies PostgreSQL migrations embedded by GoEP core and
// the addon registry generated from module manifests.
package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	coremigrations "github.com/GoEnterpricePlatform/goEP-core/internal/migrations/core"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type moduleMigrations struct {
	name    string
	fs      fs.FS
	depends []string
}

// PostgreSQLModules returns addon owners with migrations, ordered by their
// manifest dependencies.
func PostgreSQLModules() ([]string, error) {
	addons, err := orderedAddonMigrations()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(addons))
	for i, addon := range addons {
		names[i] = addon.name
	}
	return names, nil
}

// Up applies core migrations first, followed by addon migrations in
// dependency order.
func Up(databaseURL string) error {
	addons, err := orderedAddonMigrations()
	if err != nil {
		return err
	}
	if err := runModule(databaseURL, moduleMigrations{name: "core", fs: coremigrations.FS}, false); err != nil {
		return err
	}
	for _, addon := range addons {
		if err := runModule(databaseURL, addon, false); err != nil {
			return err
		}
	}
	return nil
}

// Down rolls back one migration per addon in reverse dependency order, then
// rolls back one core migration.
func Down(databaseURL string) error {
	addons, err := orderedAddonMigrations()
	if err != nil {
		return err
	}
	for i := len(addons) - 1; i >= 0; i-- {
		if err := runModule(databaseURL, addons[i], true); err != nil {
			return err
		}
	}
	return runModule(databaseURL, moduleMigrations{name: "core", fs: coremigrations.FS}, true)
}

func orderedAddonMigrations() ([]moduleMigrations, error) {
	byName := make(map[string]moduleMigrations, len(addonMigrations))
	for _, addon := range addonMigrations {
		if _, exists := byName[addon.name]; exists {
			return nil, fmt.Errorf("duplicate addon migration registry entry %q", addon.name)
		}
		addon.depends = append([]string(nil), addon.depends...)
		byName[addon.name] = addon
	}

	indegree := make(map[string]int, len(byName))
	dependents := make(map[string][]string, len(byName))
	for name, addon := range byName {
		for _, dependency := range addon.depends {
			if _, hasMigrations := byName[dependency]; !hasMigrations {
				continue
			}
			indegree[name]++
			dependents[dependency] = append(dependents[dependency], name)
		}
	}

	ready := make([]string, 0, len(byName))
	for name := range byName {
		if indegree[name] == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)
	ordered := make([]moduleMigrations, 0, len(byName))
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		ordered = append(ordered, byName[name])
		for _, dependent := range dependents[name] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.Strings(ready)
	}
	if len(ordered) != len(byName) {
		cycle := make([]string, 0)
		for name, count := range indegree {
			if count > 0 {
				cycle = append(cycle, name)
			}
		}
		sort.Strings(cycle)
		return nil, fmt.Errorf("addon migration dependency cycle detected among %v", cycle)
	}
	return ordered, nil
}

func runModule(databaseURL string, module moduleMigrations, down bool) error {
	source, err := iofs.New(module.fs, ".")
	if err != nil {
		return fmt.Errorf("open %s migration source: %w", module.name, err)
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open PostgreSQL migration connection: %w", err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "schema_migrations_" + module.name})
	if err != nil {
		return fmt.Errorf("initialize %s migration driver: %w", module.name, err)
	}
	instance, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialize %s migrations: %w", module.name, err)
	}
	defer instance.Close()
	if down {
		err = instance.Steps(-1)
	} else {
		err = instance.Up()
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply %s migrations: %w", module.name, err)
	}
	return nil
}

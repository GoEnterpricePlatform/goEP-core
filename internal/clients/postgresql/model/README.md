# PostgreSQL schema

The schema is owned by the core and its addons:

- `core.sql` contains shared platform tables.
- `addons/<module>/database/postgresql/schema.sql` contains that addon's tables.

Run `make schema` (or `go run ./cmd/goep schema`) from the repository root to
regenerate the full `schema.sql` bundle here. The bundle is a reference and
development reset artifact; the application no longer applies it at startup.

`dev_reset.sql` is generated as an optional destructive reset for a disposable
development database; it is not run by the application. Production and local
database changes use versioned migrations in each owner's migration directory
(`addons/<module>/database/postgresql/migrations/` for addons and
`internal/migrations/core/` for GoEP core). Add a new numbered `.up.sql` and
matching `.down.sql` for every schema change; do not edit migrations that have
already shipped. Each owner has its own version table
(`schema_migrations_<owner>`).

`goep generate-module-registry` generates each addon migration embed and the
central migration registry from addon manifests and migration directories for
backend addons that declare PostgreSQL support. Run
`go generate ./internal/modules/backend` after adding a module migration; the
runner then applies the generated registry in `depends` order.

Run `make migrate-up` on a new local database and before deploying a new
application version. It applies core migrations, then addon migrations in
dependency order.
`make migrate-down` rolls back one migration per module in reverse dependency
order; review the down SQL before running it because it may remove data.

## Schema design reminders

- Decide which tables need `created_at` and `updated_at`.
- Choose deletion behavior deliberately: hard delete, soft delete with
  `deleted_at`, or explicit transactional cleanup.

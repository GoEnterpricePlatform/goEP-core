# PostgreSQL schema

The Lite schema is owned by its addons:

- `addons/<module>/database/postgresql/schema.sql` contains that addon's tables.

The module lifecycle table belongs to the Managed flavor. Its source and
migration live under `internal/modules/backend/managed/`; `make schema` writes
it separately to `managed_schema.sql`.

Run `make schema` (or `go run ./cmd/goep schema`) from the repository root to
regenerate the Lite `schema.sql` bundle here. The bundle is a reference and
development reset artifact; the application no longer applies it at startup.

`dev_reset.sql` is generated as an optional destructive reset for a disposable
development database; it is not run by the application. Production and local
database changes use versioned migrations in each owner's migration directory
(`addons/<module>/database/postgresql/migrations/` for addons and
`internal/modules/backend/managed/migrations/` for lifecycle). Add a new numbered `.up.sql` and
matching `.down.sql` for every schema change; do not edit migrations that have
already shipped. Each owner has its own version table
(`schema_migrations_<owner>`).

`goep generate-module-registry` generates each addon migration embed and the
central migration registry from addon manifests and migration directories for
backend addons that declare PostgreSQL support. Run
`go generate ./internal/modules/backend` after adding a module migration; the
runner then applies the generated registry in `depends` order.

Run `make migrate-up` on a new database and before deploying a new application
version. It matches the API's `FlavorManaged`, creating the lifecycle table and
applying addon migrations. Use `make migrate-lite-up` for Lite-only databases.
`make migrate-down` rolls back the Managed migration set; `make migrate-lite-down`
rolls back the Lite set. Review down SQL before running it because it may remove
data.

## Schema design reminders

- Decide which tables need `created_at` and `updated_at`.
- Choose deletion behavior deliberately: hard delete, soft delete with
  `deleted_at`, or explicit transactional cleanup.

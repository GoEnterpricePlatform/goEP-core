# PostgreSQL schema

The schema is owned by the core and its addons:

- `core.sql` contains shared platform tables.
- `addons/<module>/database/postgresql/schema.sql` contains that addon's tables.

Run `make schema` (or `go run ./cmd/goep schema`) from the repository root
to regenerate the full `schema.sql` bundle here. It uses `IF NOT EXISTS` for
tables and indexes so development startup can apply it repeatedly without
dropping existing data. This only creates missing objects; it does not alter
existing tables. The assembler uses an explicit order: core, identity, posts,
catalog, billing. Keep this order in sync with foreign-key dependencies when
adding modules or tables.

`dev_reset.sql` is generated as an optional destructive reset for a disposable
development database; it is not run by the application. Production bootstrap
can use `schema.sql` on a new database. This is not a versioned migration
system; add that when deployed databases need safe incremental schema changes.

## Schema design reminders

- Decide which tables need `created_at` and `updated_at`.
- Choose deletion behavior deliberately: hard delete, soft delete with
  `deleted_at`, or explicit transactional cleanup.

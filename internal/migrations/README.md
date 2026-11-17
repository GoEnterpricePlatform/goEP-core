# Migraciones PostgreSQL

## Qué hace `runner.go`

`runner.go` coordina `golang-migrate`; no contiene SQL ni concatena los archivos.

- `Up` aplica primero las migraciones de core y después las de los addons, en
  orden de dependencia.
- `Down` revierte una migración por addon en orden inverso y después intenta
  revertir una de core. El SQL `down` puede borrar datos.
- Cada propietario mantiene su historial en una tabla propia, por ejemplo
  `schema_migrations_identity`.

`addons_gen.go` es generado por `go generate ./internal/modules/backend`. Enumera
los addons con migraciones PostgreSQL, embebidas por módulo, y sus dependencias.
El runner usa ese registro para decidir qué ejecutar y en qué orden. Por eso,
`make migrate-up` no descubre archivos nuevos directamente del disco.

## Al crear un módulo

`go run ./cmd/goep create-module <nombre>` crea el esqueleto y ejecuta
`go generate ./internal/modules/backend`. Si el módulo todavía no tiene
migraciones, no se agrega al registro de migraciones y no hace falta hacer nada
más.

Cuando agregues migraciones PostgreSQL al módulo:

1. Declara `postgresql` en `database_providers` del manifiesto.
2. Agrega los archivos numerados `.up.sql` y `.down.sql` en
   `addons/<nombre>/database/postgresql/migrations/`.
3. Regenera los embeds y el registro:

   ```sh
   go generate ./internal/modules/backend
   ```

4. Aplica los cambios a la base:

   ```sh
   make migrate-up
   ```

`go run ./cmd/goep schema` es independiente: genera el bundle `schema.sql` y
`dev_reset.sql`; no participa en este flujo ni aplica cambios a la base.

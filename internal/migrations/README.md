# Migraciones PostgreSQL

## Qué hace `runner.go`

`runner.go` coordina `golang-migrate`; no contiene SQL ni concatena los archivos.

- `Up` aplica las migraciones de los addons en orden de dependencia (sabor Lite).
- `UpManaged` aplica primero la migración de ciclo de vida Managed y luego las
  migraciones de addons.
- `Down` y `DownManaged` revierten sus conjuntos correspondientes. El SQL `down`
  puede borrar datos.
- Cada propietario mantiene su historial en una tabla propia, por ejemplo
  `schema_migrations_identity`.

`addons_gen.go` es generado por `go generate ./internal/modules/backend`. Enumera
los addons con migraciones PostgreSQL, embebidas por módulo, y sus dependencias.
El runner usa ese registro para decidir qué ejecutar y en qué orden. Por eso,
`make migrate-up` no descubre archivos nuevos directamente del disco. Ahora
ejecuta Managed, que crea `goep_core_modules`, para coincidir con el flavor
configurado en la API. `make migrate-lite-up` ejecuta solo las migraciones de
addons. El store PostgreSQL valida que exista `goep_core_modules` al iniciar;
la tabla no se crea desde el código de arranque.

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

`go run ./cmd/goep schema` es independiente: genera el bundle Lite `schema.sql`,
`dev_reset.sql` y el `managed_schema.sql` opcional; no aplica cambios a la base.

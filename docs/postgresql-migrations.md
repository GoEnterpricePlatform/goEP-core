# Esquemas y migraciones PostgreSQL

Esta guía distingue el ensamblado del esquema (`goep schema`) de las
migraciones versionadas que cambian una base de datos.

## Qué hace `go run ./cmd/goep schema`

Desde la raíz del repositorio, el comando lee estos archivos fuente:

| Propietario | Archivo fuente |
| --- | --- |
| GoEP core | `internal/clients/postgresql/model/core.sql` |
| identity | `addons/identity/database/postgresql/schema.sql` |
| posts | `addons/posts/database/postgresql/schema.sql` |
| catalog | `addons/catalog/database/postgresql/schema.sql` |
| billing | `addons/billing/database/postgresql/schema.sql` |

Con ellos genera dos archivos bajo `internal/clients/postgresql/model/`:

- `schema.sql`: una copia ensamblada de los esquemas fuente, con `IF NOT EXISTS`
  en tablas e índices.
- `dev_reset.sql`: primero elimina las tablas incluidas y luego agrega el
  esquema ensamblado. **Es destructivo; úsalo solo en una base descartable de
  desarrollo.**

El comando no se conecta a PostgreSQL, no ejecuta SQL y no cambia el historial
de migraciones. La aplicación tampoco ejecuta `schema.sql` al arrancar. Estos
archivos sirven como referencia y para el flujo de reset local.

Actualmente la lista de archivos fuente está escrita explícitamente en
`internal/cli/schema.go`. Por eso, crear un addon nuevo no lo agrega
automáticamente a `schema.sql`.

## Cuándo ejecutar `schema`

Ejecuta `go run ./cmd/goep schema` —o `make schema`— después de modificar uno de
los archivos fuente `schema.sql` y quieras actualizar el bundle de referencia y
el reset local.

No es necesario ejecutarlo al crear un módulo vacío. `go run
./cmd/goep create-module <nombre>` crea el esqueleto del addon y regenera los
registros Go, pero no crea tablas ni ejecuta `schema`.

Si el addon nuevo va a usar PostgreSQL y quieres que su esquema aparezca en el
bundle generado, crea `addons/<nombre>/database/postgresql/schema.sql`, añade
ese archivo a la lista de `files` en `internal/cli/schema.go` y luego ejecuta
`go run ./cmd/goep schema`. El orden debe respetar las claves foráneas entre
módulos.

## Cómo agregar una migración

Las migraciones aplicables a una base viven con el propietario del esquema:

```text
addons/<nombre>/database/postgresql/migrations/
├── 000001_create_example.up.sql
└── 000001_create_example.down.sql
```

GoEP core mantiene las suyas en `internal/migrations/core/`. Usa números
crecientes por módulo. Cada `up` describe el cambio hacia adelante y cada
`down` su reversión. Una migración que ya se desplegó se considera inmutable:
para modificar el esquema después, agrega otra migración.

Para un addon:

1. Declara `"postgresql"` en `database_providers` dentro de `manifest.json`.
2. Crea el par numerado `.up.sql` y `.down.sql` en la carpeta del módulo.
3. Regenera el código de embed y el registro desde la raíz:

   ```powershell
   go generate ./internal/modules/backend
   ```

   `create-module` ya ejecuta este paso al crear el esqueleto. Si agregas las
   migraciones después, vuelve a ejecutarlo.

4. Revisa `addons/<nombre>/database/postgresql/migrations/migrations_gen.go` y
   `internal/migrations/addons_gen.go`; son archivos generados y no se editan
   manualmente.
5. Aplica las migraciones a la base deseada:

   ```powershell
   make migrate-up
   ```

El registro generado incluye los addons backend que declaran PostgreSQL y
tienen migraciones `.up.sql`. El runner los ordena por `depends`. Los addons
web y los que no declaran PostgreSQL no se incluyen.

## Flujo al crear un módulo con el CLI

```text
go run ./cmd/goep create-module reports
    └─ crea addons/reports/ y su manifest.json
    └─ ejecuta go generate ./internal/modules/backend
       └─ regenera fábricas, assets y el registro de migraciones
```

El esqueleto no contiene un esquema PostgreSQL ni migraciones, así que el
registro de migraciones no tendrá una entrada para el módulo hasta que añadas
un archivo `.up.sql` y regeneres. El comando `schema` es independiente y no se
ejecuta durante `create-module`.

## Aplicar y revertir en bases de datos

`make migrate-up` aplica primero las migraciones de core y luego las de los
addons, en orden de dependencias. Cada propietario mantiene su propia tabla de
versiones: `schema_migrations_core` y `schema_migrations_<módulo>`.

`make migrate-down` baja una migración por addon en orden inverso y después
intenta bajar una de core. Revisa el SQL `down` antes de usarlo: puede eliminar
tablas o datos. En producción, el flujo normal es respaldar la base, aplicar
`migrate-up` antes de iniciar la nueva versión y observar el resultado; no
ejecutar `down` automáticamente como procedimiento de despliegue.

La API no aplica migraciones al arrancar. En desarrollo y producción, ejecuta
el comando de migraciones explícitamente contra la base que corresponda.

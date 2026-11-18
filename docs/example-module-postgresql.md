# Ejemplo: módulo PostgreSQL con tablas

Este ejemplo crea el módulo `inventario` con una tabla de productos.

## 1. Crear el módulo

Desde la raíz del repositorio:

```sh
go run ./cmd/goep create-module inventario
```

El CLI pregunta qué base de datos usará el módulo. PostgreSQL es la primera
opción y la predeterminada si presionas Enter. También puedes elegir MongoDB,
ambas bases o ninguna. Para evitar la pregunta, indica `--database`, por
ejemplo `--database postgresql`.

Al elegir PostgreSQL (o ambas bases), el CLI crea esta base de migración:

```text
addons/inventario/database/postgresql/migrations/
├── 000001_initial.up.sql
└── 000001_initial.down.sql
```

Los archivos incluyen comentarios para completar con las tablas iniciales y
su reversión. El CLI regenera los embeds y registros Go al terminar.

## 2. Declarar PostgreSQL

El manifiesto ya declara PostgreSQL cuando se selecciona esa opción:

```json
{
  "name": "inventario",
  "version": "1.0.0",
  "summary": "Gestión de inventario",
  "depends": [],
  "database_providers": ["postgresql"]
}
```

Si el módulo también usa MongoDB, incluye ambos proveedores en la lista.

## 3. Completar la migración inicial

Edita `000001_initial.up.sql` y agrega la tabla de productos:

```sql
CREATE TABLE tb_inventory_product (
    id         UUID PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    stock      INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

En `000001_initial.down.sql`, escribe la reversión:

```sql
DROP TABLE IF EXISTS tb_inventory_product;
```

Las migraciones ya aplicadas no se editan. Para cambiar la tabla después,
agrega `000002_...up.sql` y su reversión.

## 4. Regenerar el registro

Después de agregar los archivos SQL, ejecuta:

```sh
go generate ./internal/modules/backend
```

Esto genera el embed del módulo en su carpeta `migrations/` y lo agrega a
`internal/modules/migrations/addons_gen.go`. El CLI de creación no puede hacerlo por
adelantado porque todavía no existían las migraciones.

## 5. Aplicar la migración

Configura `PG_DB_URI` para la base destino y ejecuta:

```sh
make migrate-up
```

El runner aplicará primero las migraciones de core y luego las de `inventario`
según sus dependencias. PostgreSQL registrará el avance en
`schema_migrations_inventario`.

`go run ./cmd/goep schema` no es necesario para este flujo: solo regenera el
bundle de referencia `schema.sql` y el reset de desarrollo; no ejecuta SQL en
la base.

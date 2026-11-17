# Ejemplo: módulo PostgreSQL con tablas

Este ejemplo crea el módulo `inventario` con una tabla de productos.

## 1. Crear el módulo

Desde la raíz del repositorio:

```sh
go run ./cmd/goep create-module inventario
```

El CLI crea `addons/inventario/` y regenera los registros Go. Todavía no hay
migraciones para este módulo.

## 2. Declarar PostgreSQL

En `addons/inventario/manifest.json`, agrega el proveedor PostgreSQL:

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

## 3. Crear la migración inicial

Crea el directorio `addons/inventario/database/postgresql/migrations/` y el
archivo `000001_create_products.up.sql`:

```sql
CREATE TABLE tb_inventory_product (
    id         UUID PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    stock      INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

En el archivo correspondiente `000001_create_products.down.sql`, escribe la
reversión:

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
`internal/migrations/addons_gen.go`. El CLI de creación no puede hacerlo por
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

# Documentación de arquitectura modular

La guía general está en [Arquitectura modular de GoEP Core](modular-architecture.md).

Para agregar un módulo paso a paso, consulta [Agregar un módulo de backend](add-module.md).

Para distinguir los esquemas fuente del ensamblado y las migraciones PostgreSQL,
consulta [Esquemas y migraciones PostgreSQL](postgresql-migrations.md).

Para ver el proceso completo con SQL de ejemplo, consulta
[Ejemplo de módulo PostgreSQL](example-module-postgresql.md).

Para entender el runner y el flujo de ejecución por módulo, consulta
[Migraciones PostgreSQL](../internal/migrations/README.md).

## CLI del repositorio

Desde la raíz del repositorio, ejecuta el CLI sin argumentos para abrir el menú
interactivo:

```powershell
go run ./cmd/goep
```

El menú permite crear un módulo o ensamblar el esquema PostgreSQL. También
puedes invocar las acciones directamente desde la consola.

### Crear el esqueleto de un módulo

Desde la raíz del repositorio, ejecuta el CLI e indica el nombre del módulo en minúsculas. Usa guiones para separar palabras:

```powershell
go run ./cmd/goep create-module invoices
```

También puedes definir el resumen que quedará en el manifiesto:

```powershell
go run ./cmd/goep create-module --summary "Sales reports" sales-reports
```

El comando crea `manifest.json`, los paquetes base `config/` y `domain/`, y `static/` con `assets.json`, `src/ts/`, `src/css/` y `public/` bajo `addons/<nombre>/`. Luego ejecuta `go generate ./internal/modules/backend` automáticamente para agregarlo al registro compilado. Si ya existe una carpeta con ese nombre, el CLI se detiene sin sobrescribirla. Después completa el constructor con los servicios, rutas y assets que necesite el módulo; el esqueleto no agrega lógica de negocio ni compila assets todavía.

## Documentos de referencia

- [Composición e inyección de servicios](backend-module-composition.md)
- [Descubrimiento de manifiestos](module-discovery.md)
- [Registro compilado y dependencias](compiled-module-registry.md)
- [Guía detallada del generador del registro](generate-module-registry.md)
- [Hoja de ruta para assets de Web App por módulo](web-app-module-assets.md)
- [Persistencia del ciclo de vida](module-lifecycle.md)

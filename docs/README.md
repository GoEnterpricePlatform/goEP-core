# Documentación de arquitectura modular

La guía general está en [Arquitectura modular de GoEP Core](modular-architecture.md).

Para agregar un módulo paso a paso, consulta [Agregar un módulo de backend](add-module.md).

Para distinguir los esquemas fuente del ensamblado y las migraciones PostgreSQL,
consulta [Esquemas y migraciones PostgreSQL](postgresql-migrations.md).

Para ver el proceso completo con SQL de ejemplo, consulta
[Ejemplo de módulo PostgreSQL](example-module-postgresql.md).

Para entender el runner y el flujo de ejecución por módulo, consulta
[Migraciones PostgreSQL](postgresql-migrations.md#aplicar-y-revertir-en-bases-de-datos).

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
go run ./cmd/goep create-module --summary "Sales reports" --ui dashboard --ui-layout replace sales-reports
```

El comando pregunta qué base de datos, qué tipo de interfaz y si se debe mostrar la barra `GOEP CORE`. `--ui` acepta `simple` o `dashboard`; dashboard siempre usa el sidebar de GoEP. `--ui-layout` acepta `topbar` (incluye la barra `GOEP CORE`) o `replace` (la oculta). Si omites `--ui-layout`, el CLI selecciona `replace`. Las tarjetas abren la UI generada en `/v1/<nombre>/`.

El comando crea `manifest.json`, los paquetes base `config/` y `domain/`, `ui/handler/` y `ui/pages/`, y `static/` con `assets.json`, `src/ts/`, `src/css/` y `public/` bajo `addons/<nombre>/`. Después genera las plantillas Templ y ejecuta `go generate ./internal/modules/backend` para agregar el módulo al registro compilado. Si ya existe una carpeta con ese nombre, el CLI se detiene sin sobrescribirla. Completa el constructor con los servicios, rutas y assets que necesite el módulo; el esqueleto no agrega lógica de negocio ni compila assets todavía.

Para `--ui dashboard`, el CLI también crea `ui/navigation/menu.go` con un enlace plano y un grupo expandible de ejemplo. Edita `AdminMenu` allí para agregar enlaces o grupos con `Children`. El handler generado resuelve la ruta hija y actualiza el título/estado activo; para agregar contenido específico, registra su ruta en `ui/handler/` y crea la vista Templ dentro de `ui/pages/`. El ejemplo `addons/billing` organiza una vista por proveedor en `ui/pages/payments/` y `ui/pages/subs/`.

## Documentos de referencia

- [Composición e inyección de servicios](backend-module-composition.md)
- [Descubrimiento de manifiestos](module-discovery.md)
- [Registro compilado y dependencias](compiled-module-registry.md)
- [Guía detallada del generador del registro](generate-module-registry.md)
- [Hoja de ruta para assets de Web App por módulo](web-app-module-assets.md)
- [Persistencia del ciclo de vida](module-lifecycle.md)

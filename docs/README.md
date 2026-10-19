# Documentación de arquitectura modular

La guía general está en [Arquitectura modular de GoEP Core](modular-architecture.md).

Para agregar un módulo paso a paso, consulta [Agregar un módulo de backend](add-module.md).

## Crear el esqueleto desde la consola

Desde la raíz del repositorio, ejecuta el CLI e indica el nombre del módulo en minúsculas. Usa guiones para separar palabras:

```powershell
go run ./cmd/create-module invoices
```

También puedes definir el resumen que quedará en el manifiesto:

```powershell
go run ./cmd/create-module --summary "Sales reports" sales-reports
```

El comando crea `manifest.json` y los paquetes base `config/` y `domain/` bajo `addons/<nombre>/`. Luego ejecuta `go generate ./internal/modules/backend` automáticamente para agregarlo al registro compilado. Si ya existe una carpeta con ese nombre, el CLI se detiene sin sobrescribirla. Después completa el constructor con los servicios y rutas que necesite el módulo; el esqueleto no agrega lógica de negocio.

## Documentos de referencia

- [Composición e inyección de servicios](backend-module-composition.md)
- [Descubrimiento de manifiestos](module-discovery.md)
- [Registro compilado y dependencias](compiled-module-registry.md)
- [Guía detallada del generador del registro](generate-module-registry.md)
- [Persistencia del ciclo de vida](module-lifecycle.md)

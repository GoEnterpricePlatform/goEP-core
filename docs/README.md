# Documentación de arquitectura modular

La guía general está en [Arquitectura modular de GoEP Core](modular-architecture.md).

Para agregar un módulo paso a paso, consulta [Agregar un módulo de backend](add-module.md).

## Generador del registro de módulos

`cmd/generate-module-registry/main.go` es una herramienta de desarrollo. Lee los manifiestos de `addons/` y genera `internal/modules/backend/factories_gen.go`, que importa y registra los módulos disponibles para compilarlos dentro del ejecutable.

El generador **no se ejecuta cuando arranca la API**, ni automáticamente con `go build` o `go test`. Se ejecuta mediante la directiva `//go:generate` de `internal/modules/backend/factories.go`. Después de agregar o quitar un módulo de `addons/`, corre desde la raíz del repositorio:

```powershell
go generate ./internal/modules/backend
```

Ese comando invoca internamente `go run ../../../cmd/generate-module-registry` y actualiza `factories_gen.go`. Revisa e incluye ese archivo generado en el mismo cambio de código. No lo edites manualmente. Si no ejecutas el generador, el nuevo módulo no queda importado en el binario y la validación del registro fallará al descubrir su manifiesto.

No hace falta ejecutarlo para cambios internos de un módulo que no alteren el manifiesto ni la lista de módulos.

## Documentos de referencia

- [Composición e inyección de servicios](backend-module-composition.md)
- [Descubrimiento de manifiestos](module-discovery.md)
- [Registro compilado y dependencias](compiled-module-registry.md)
- [Persistencia del ciclo de vida](module-lifecycle.md)

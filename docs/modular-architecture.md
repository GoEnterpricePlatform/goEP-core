# Arquitectura modular de GoEP Core

Esta guía describe cómo agregar, compilar, instalar y habilitar módulos de backend. Los módulos oficiales se enlazan estáticamente: el manifiesto permite descubrir metadatos, pero no carga código Go desde disco. Agregar una implementación requiere registrarla y compilar una nueva versión del ejecutable.

## Estados de un módulo

- **Disponible:** su manifiesto fue descubierto y existe una fábrica compilada para el mismo nombre.
- **Instalado:** el estado persistido indica `installed: true`.
- **Habilitado:** el estado indica `enabled: true`; sus rutas API se registran en el mux público durante el arranque.

El estado es global para la instalación de GoEP Core y se guarda en MongoDB, en `goep_core_modules`. Una instalación con esa colección vacía instala y habilita los módulos oficiales descubiertos para conservar el comportamiento actual. En arranques posteriores se utiliza el estado guardado. Si la colección no está vacía, los módulos sin documento quedan deshabilitados.

El cambio de estado se aplica en el siguiente arranque. Web App todavía no forma parte del ciclo de vida: los servicios backend se construyen para que sus páginas sigan funcionando, mientras que las rutas API de módulos deshabilitados quedan aisladas.

## Agregar un módulo backend

1. Crea el directorio y el manifiesto bajo `addons/`. El nombre debe coincidir con el nombre de la fábrica compilada:

   ```json
   {
     "name": "reports",
     "version": "1.0.0",
     "summary": "Report generation",
     "depends": ["identity"]
   }
   ```

   Para el módulo de IA existente, el directorio y el nombre son `ai-tool-calling`.

2. Implementa el módulo y su constructor Go con configuración específica. Declara servicios requeridos en un `ModuleDeps` tipado y consume contratos de puerto, sin importar la implementación privada de otro módulo.

3. Registra la fábrica en `internal/modules/backend/factories.go`. El nombre debe ser idéntico al `name` del manifiesto. La fábrica vincula el constructor al ejecutable.

4. Añade la construcción y la inyección de servicios en `internal/modules/backend/composition.go`, respetando el orden de dependencias. El resolver rechaza manifiestos sin fábrica, fábricas sin manifiesto, dependencias faltantes y ciclos antes de inicializar los módulos.

5. Añade pruebas para el manifiesto, la fábrica, las dependencias, los contratos requeridos y las rutas públicas del módulo.

6. Ejecuta las pruebas del sistema modular y genera un nuevo binario:

   ```powershell
   go test ./internal/modules/...
   go test ./cmd/api/...
   go build ./cmd/api/env/dev/main
   go build ./cmd/api/env/prod
   ```

## Administrar el ciclo de vida desde Go

El bootstrap devuelve el `*lifecycle.Manager` en `backend.Modules.Lifecycle`. El manager ofrece `List`, `Install`, `Enable` y `Disable`; comprueba que las dependencias estén instaladas y habilitadas, y no elimina datos del módulo al deshabilitarlo.

Ejemplo desde código de backend que ya recibió `backend.Modules`:

```go
if err := backend.Lifecycle.Install(ctx, "reports"); err != nil {
    return err
}
if err := backend.Lifecycle.Enable(ctx, "reports"); err != nil {
    return err
}
```

Para deshabilitarlo:

```go
if err := backend.Lifecycle.Disable(ctx, "reports"); err != nil {
    return err
}
```

La aplicación aún no expone endpoints ni un comando CLI para administrar módulos. Las operaciones anteriores son la API interna para integrarlas desde una interfaz administrativa. Evita editar documentos directamente: el manager valida dependencias y mantiene las marcas de tiempo y versión.

La desinstalación no está implementada. La política actual es deshabilitar primero y, cuando se implemente la operación, retirar solo el documento de estado; las colecciones de datos no se borrarán automáticamente.

## Pruebas con MongoDB

La suite normal no requiere un servidor MongoDB. Para correr además la prueba de persistencia real entre dos conexiones, configura `GOEP_TEST_MONGODB_URI` apuntando a una instancia de pruebas:

```powershell
$env:GOEP_TEST_MONGODB_URI = "mongodb://localhost:27017"
go test ./internal/modules/lifecycle -run TestMongoStorePersistsStateAcrossClients
```

La prueba usa una base temporal con nombre único y la elimina al finalizar. No apuntes esta variable a una base con datos importantes.

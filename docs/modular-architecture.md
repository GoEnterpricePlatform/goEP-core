# Arquitectura modular de GoEP Core

Esta guía describe cómo agregar y compilar addons de backend. Los módulos oficiales se enlazan estáticamente: el manifiesto permite descubrir metadatos, pero no carga código Go desde disco. Agregar una implementación requiere regenerar el registro y compilar una nueva versión del ejecutable.

## Sabores de composición

`backend.Config.Flavor` selecciona cómo se componen los addons. Si se omite, el valor vacío usa `backend.FlavorLite`. El servidor principal (`cmd/api/server/v1`) declara `FlavorManaged` explícitamente.

- **Lite (`lite`, predeterminado):** descubre y valida los addons, resuelve sus dependencias y los inicializa en orden. Todos los módulos backend compatibles con el proveedor de base de datos quedan cargados y sus rutas API se registran. No consulta ni crea estado de ciclo de vida y no usa la colección `goep_core_modules`.
- **Managed (`managed`):** conserva el ciclo de vida existente de instalación y habilitación persistidas. Es el sabor anterior, disponible para retomarlo cuando se implemente la administración de módulos desde la aplicación.

El registro generado sigue enlazando los addons estáticamente en ambos sabores. "Lite" simplifica su carga en tiempo de arranque; no instala código Go desde el directorio de addons.

La implementación está separada en `internal/modules/backend/lite/` y `internal/modules/backend/managed/`. El despacho, los tipos de configuración y la composición de servicios compartidos permanecen en `internal/modules/backend/`. Las migraciones y el esquema PostgreSQL de `goep_core_modules` también pertenecen a `managed/`; `make migrate-up` aplica Managed para coincidir con el servidor actual, mientras que `make migrate-lite-up` aplica solo las migraciones de addons.

## Ciclo de vida del sabor managed

Esta sección documenta el sabor `managed`; no describe el arranque predeterminado.

- **Disponible:** su manifiesto fue descubierto y existe una fábrica compilada para el mismo nombre.
- **Instalado:** el estado persistido indica `installed: true`.
- **Habilitado:** el estado indica `enabled: true`; sus rutas API se registran en el mux público durante el arranque.

El estado es global para la instalación de GoEP Core y se guarda en MongoDB, en `goep_core_modules`. Una instalación con esa colección vacía instala y habilita los módulos oficiales descubiertos. En arranques posteriores se utiliza el estado guardado. Si la colección no está vacía, los módulos sin documento quedan deshabilitados.

El cambio de estado se aplica en el siguiente arranque. Web App todavía no forma parte del ciclo de vida: los servicios backend se construyen para que sus páginas sigan funcionando, mientras que las rutas API de módulos deshabilitados quedan aisladas.

## Agregar un módulo backend

1. Crea el directorio y el manifiesto bajo `addons/`:

   ```json
   {
     "name": "reports",
     "version": "1.0.0",
     "summary": "Report generation",
     "depends": ["identity"]
   }
   ```

   Para el módulo de IA existente, el directorio y el nombre son `ai-tool-calling`.

2. Implementa el módulo y su constructor Go en `addons/<nombre>/config/`. El paquete debe exponer exactamente un constructor con esta forma: `func New...Module(ModuleConfig) (*Module, error)`. El generador lo detecta automáticamente.

3. Implementa `RegisterBackendModule(*modules.Registry) error` en el mismo paquete. Este adaptador arma el `ModuleConfig` tipado, obtiene dependencias con `modules.RequireService[T]` y publica contratos para otros módulos con `modules.ProvideService[T]`. Declara en `depends` los módulos que proveen esos contratos.

4. Regenera el registro compilado:

   ```powershell
   go generate ./internal/modules/backend
   ```

   El generador lee los manifiestos directos de `addons/` y actualiza `internal/modules/backend/factories_gen.go`. No edites ese archivo generado ni agregues imports manuales a `factories.go`.

5. El backend inicializa el adaptador automáticamente según `depends`; no agregues una llamada al constructor en `internal/modules/backend/composition.go`. Si otro código de infraestructura, como `cmd/api/server/v1`, necesita un servicio del nuevo módulo, expón ese módulo desde `backend.Modules` y conecta ese servicio allí.

6. Añade pruebas para el manifiesto, la fábrica, las dependencias, los contratos requeridos y las rutas públicas del módulo.

7. En la última parte, ejecuta las pruebas del sistema modular y genera un nuevo binario:

   ```powershell
   go generate ./internal/modules/backend
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

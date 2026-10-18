# Ejemplo: agregar un módulo de backend

Esta guía muestra qué archivos crear y cuáles modificar para agregar un módulo GoEP Core llamado `reports`. El ejemplo sigue la composición actual de `posts`, `billing` y `catalog`: manifiestos para metadatos, fábricas compiladas y configuración tipada.

> Este README es una guía. No crees `addons/reports/` hasta que también vayas a registrar su fábrica y su composición: el descubridor requiere que cada subdirectorio de `addons/` tenga un `manifest.json`, y el bootstrap rechaza manifiestos sin implementación compilada.

## 1. Estructura de archivos

Como punto de partida, crea solo las carpetas que el módulo necesita:

```text
addons/reports/
├── manifest.json                 # Obligatorio
├── config/
│   └── module_config.go          # Obligatorio: configuración y constructor
├── domain/                       # Si hay entidades y reglas de dominio
├── port/
│   └── ports.go                  # Si otro paquete consume sus servicios
├── repository/
│   └── mongo/                    # Si persiste datos en MongoDB
├── service/                      # Reglas de aplicación
├── handler/                      # Si expone rutas HTTP
├── model/                        # Si necesita mapear documentos Mongo
└── *_test.go                     # Pruebas del módulo
```

No todos los módulos necesitan `domain`, `port`, `model` o `repository`. Mantén los puertos dentro de `port/` cuando otra capa o módulo necesite depender de una interfaz; evita crear paquetes vacíos.

## 2. Crear el manifiesto

Crea `addons/reports/manifest.json`:

```json
{
  "name": "reports",
  "version": "1.0.0",
  "summary": "Report generation",
  "depends": ["identity"]
}
```

Reglas para este archivo:

- `name` y `version` son obligatorios.
- `name` debe coincidir con el nombre registrado en la fábrica.
- `depends` contiene nombres de módulos requeridos, no nombres de servicios o paquetes Go.
- Declara `identity` si el módulo necesita autenticación o depende funcionalmente de Identity. Si solo necesita infraestructura compartida, no agregues esa dependencia por costumbre.
- Un directorio directo de `addons/` sin `manifest.json` hace fallar el descubrimiento.

Los nombres actuales son `identity`, `posts`, `billing`, `catalog` y `ai-tool-calling`. El directorio de AI también se llama `addons/ai-tool-calling/`.

## 3. Crear la configuración y el constructor

El compositor entrega dependencias concretas mediante una estructura propia del módulo. Por ejemplo, si Reports expone rutas protegidas, puede recibir el contrato de autenticación que ya usa Posts:

```go
package config

import (
    "net/http"

    middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
    appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
    AppEnvs    *appConfig.AppEnvs
    AppClients *appConfig.AppClients
    APIv1      *http.ServeMux
    DB         *mongo.Database
    Deps       ModuleDeps
}

type ModuleDeps struct {
    AuthApiMdw middlewarePort.Authenticator
}

type Module struct {
    // Exporta aquí solo los servicios que realmente consumen otros módulos
    // o Web App, preferiblemente mediante interfaces de port/.
}

func NewReportsModule(cfg ModuleConfig) (*Module, error) {
    // Construye repositorios y servicios, luego registra los handlers.
    // Propaga errores con return nil, err; no los ocultes.
    return &Module{}, nil
}
```

Este esqueleto no implementa almacenamiento ni rutas. Si Reports no usa Mongo, elimina `DB`; si no requiere autenticación, elimina `Deps` y no agregues `identity` al manifiesto por ese motivo. Compara la implementación completa de [Posts](../addons/posts/config/module_config.go) antes de copiar patrones de Mongo o handlers.

Las rutas se registran durante la construcción usando `cfg.APIv1`. Si el módulo necesita autenticación, sus handlers deben recibir `cfg.Deps.AuthApiMdw`, igual que [el handler de Posts](../addons/posts/handler/handler.go). El compositor entrega el mux público cuando el módulo está habilitado y un mux aislado cuando está deshabilitado.

## 4. Agregar contratos si otro módulo consume Reports

Si otro módulo necesita consultar reportes, define una interfaz pequeña en `addons/reports/port/ports.go`. Por ejemplo:

```go
package port

import "context"

type ReportService interface {
    Generate(ctx context.Context, reportID string) ([]byte, error)
}
```

El paquete del consumidor importa `addons/reports/port`, no `addons/reports/service` ni `addons/reports/config`. Añade el servicio a su propio `ModuleDeps` y conecta la interfaz desde el compositor. También agrega el nombre de Reports a `depends` del manifiesto consumidor para que el lifecycle manager aplique el orden y los requisitos de habilitación.

No agregues `port/` solo para duplicar interfaces que nadie consume. En Posts, por ejemplo, los contratos compartidos están en `addons/posts/port/ports.go`.

## 5. Registrar la fábrica compilada

Modifica `internal/modules/backend/factories.go`. Agrega el import:

```go
reportsModule "github.com/GoEnterpricePlatform/goEP-core/addons/reports/config"
```

Y registra el constructor dentro de `NewCompiledRegistry`:

```go
if err := modules.RegisterFactory(registry, modules.Factory[
    reportsModule.ModuleConfig,
    *reportsModule.Module,
]{
    Name: "reports",
    New:  reportsModule.NewReportsModule,
}); err != nil {
    return nil, err
}
```

El tipo de configuración y el tipo de retorno son específicos de Reports. No reemplaces estas fábricas tipadas por un contenedor genérico. El descubridor lee JSON; no carga automáticamente paquetes Go.

## 6. Conectar la composición

Modifica `internal/modules/backend/composition.go`:

1. Importa `addons/reports/config`.
2. Obtén su fábrica con `modules.GetFactory[reportsModule.ModuleConfig, *reportsModule.Module](registry, "reports")`.
3. Construye el módulo pasando `AppEnvs`, `AppClients`, `DB`, el mux que devuelve `moduleRoutes(...)` y los contratos requeridos.
4. Coloca la llamada después de los módulos de los que depende y antes de los consumidores de Reports.
5. Si Web App u otro módulo consume un servicio de Reports, agrega ese servicio a `backend.Modules` o al `ModuleDeps` tipado del consumidor y conecta el valor aquí.

Ejemplo para Reports protegido por Identity, después de crear Identity:

```go
reports, err := reportsFactory.New(reportsModule.ModuleConfig{
    AppEnvs:    cfg.AppEnvs,
    AppClients: cfg.AppClients,
    APIv1:      moduleRoutes(cfg.APIv1, enabled, "reports"),
    DB:         cfg.AppClients.DB,
    Deps: reportsModule.ModuleDeps{
        AuthApiMdw: identity.AuthApiMdw,
    },
})
if err != nil {
    return nil, err
}
```

Si `Module` exporta un servicio usado por otro consumidor, devuelve la instancia de Reports dentro de `backend.Modules` y conecta su servicio explícitamente. `depends` valida el ciclo de vida; no realiza esa inyección por sí solo.

## 7. Agregar pruebas

Incluye pruebas según lo que implemente el módulo:

- `manifest_test.go`: nombre, versión y dependencias correctas.
- Pruebas unitarias de servicios con repositorios en memoria o mocks.
- Pruebas de handlers para respuestas HTTP y autenticación.
- Pruebas del repositorio Mongo para consultas e índices.
- Pruebas del registro y composición para fábrica disponible, dependencias y rutas.

El registro oficial se contrasta con los manifiestos por `TestCompiledBackendRegistryMatchesOfficialManifests`. Si agregas Reports al registro, actualiza las expectativas de esa prueba.

Ejecuta:

```powershell
go test ./internal/modules/...
go test ./addons/reports/...
go test ./...
```

La prueba Mongo de ciclo de vida requiere `GOEP_TEST_MONGODB_URI`; consulta [la guía de arquitectura modular](modular-architecture.md#pruebas-con-mongodb).

## 8. Instalar y habilitar

Después de agregar Reports, compila un nuevo ejecutable y despliega el manifiesto junto con el código. Si la colección de lifecycle ya contiene estados, Reports se descubre pero queda deshabilitado hasta instalarlo y habilitarlo con el manager:

```go
if err := backend.Lifecycle.Install(ctx, "reports"); err != nil {
    return err
}
if err := backend.Lifecycle.Enable(ctx, "reports"); err != nil {
    return err
}
```

El manager rechaza `Enable` si Reports no está instalado o si una dependencia requerida no está instalada y habilitada. Deshabilitarlo conserva su documento de estado y sus colecciones; sus rutas dejan de exponerse en el siguiente arranque.

La aplicación aún no incluye una pantalla, endpoint HTTP ni comando CLI para estas operaciones. Usa el manager desde una integración administrativa de confianza. En una instalación completamente nueva, el bootstrap instala y habilita los módulos descubiertos una vez para conservar el comportamiento inicial.

## Archivos que se modifican

| Archivo | Cuándo modificarlo |
|---|---|
| `addons/reports/manifest.json` | Siempre; metadatos y dependencias del nuevo módulo. |
| `addons/reports/config/module_config.go` | Siempre; configuración tipada, servicios exportados y constructor. |
| `internal/modules/backend/factories.go` | Siempre; enlaza el constructor al ejecutable. |
| `internal/modules/backend/composition.go` | Siempre; orden de construcción e inyección explícita. |
| Manifiesto del consumidor | Si otro módulo depende de Reports a nivel de lifecycle. |
| `ModuleDeps` del consumidor y compositor | Si otro módulo consume un servicio de Reports. |
| `internal/modules/backend/factories_test.go` | Si cambian las expectativas de fábricas oficiales. |
| Pruebas del módulo | Para validar contratos, servicios, persistencia y rutas. |

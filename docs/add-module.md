# Ejemplo: agregar un módulo de backend

Esta guía muestra qué archivos crear y cuáles modificar para agregar un módulo GoEP Core llamado `reports`. El registro de fábricas se genera desde los manifiestos; la conexión de dependencias entre módulos sigue siendo explícita en la composición.

> Este README es una guía. El descubridor requiere que cada subdirectorio de `addons/` tenga un `manifest.json`, y el bootstrap rechaza manifiestos sin implementación compilada.

## 1. Estructura de archivos

Como punto de partida, crea solo las carpetas que el módulo necesita:

```text
addons/reports/
├── manifest.json                 # Obligatorio
├── static_assets_gen.go          # Generado por go generate; embed del módulo
├── config/
│   ├── module_config.go          # Obligatorio: configuración y constructor
│   └── backend_module.go         # Obligatorio: adaptador de registro e inicio
├── static/                       # Lo crea el CLI; assets opcionales de Web App
│   ├── assets.json               # Puntos de entrada TS/CSS
│   ├── src/ts/
│   ├── src/css/
│   └── public/                    # Imágenes y archivos públicos
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

El CLI crea `static/` con un manifiesto vacío aunque el módulo todavía no tenga interfaz web. Al terminar, el generador crea `static_assets_gen.go` en la raíz del addon y actualiza los dos registros generados del backend. No edites esos archivos a mano. Consulta el [contrato de assets de Web App](web-app-module-assets.md#contrato-de-staticassetsjson) para declarar entradas TS/CSS y organizar recursos públicos. Al ejecutar el compilador de assets, la salida se genera en `static/dist/` dentro del mismo módulo.

## 2. Crear el manifiesto

Crea `addons/reports/manifest.json`:

```json
{
  "name": "reports",
  "version": "1.0.0",
  "summary": "Report generation",
  "icon": "static/public/icon.svg",
  "depends": ["identity"]
}
```

Reglas para este archivo:

- `name` y `version` son obligatorios.
- `icon` es opcional; cuando existe, apunta a un archivo dentro de `static/public/`. La Web App lo publica bajo `/static/modules/<nombre>/public/`.
- `name` debe coincidir con el nombre registrado en la fábrica.
- `depends` contiene nombres de módulos requeridos, no nombres de servicios o paquetes Go.
- Declara `identity` si el módulo necesita autenticación o depende funcionalmente de Identity. Si solo necesita infraestructura compartida, no agregues esa dependencia por costumbre.
- Un directorio directo de `addons/` sin `manifest.json` hace fallar el descubrimiento.

Los nombres actuales son `identity`, `posts`, `billing`, `catalog` y `ai-tool-calling`. El directorio de AI también se llama `addons/ai-tool-calling/`.

## 3. Crear la configuración y el constructor

El módulo conserva una configuración tipada. Por ejemplo, si Reports expone rutas protegidas, declara el contrato de autenticación que necesita:

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

Las rutas se registran durante la construcción usando `cfg.APIv1`. Si el módulo necesita autenticación, sus handlers deben recibir `cfg.Deps.AuthApiMdw`, igual que [el handler de Posts](../addons/posts/handler/handler.go). El runtime entrega el mux público cuando el módulo está habilitado y un mux aislado cuando está deshabilitado.

## 4. Agregar contratos si otro módulo consume Reports

Si otro módulo necesita consultar reportes, define una interfaz pequeña en `addons/reports/port/ports.go`. Por ejemplo:

```go
package port

import "context"

type ReportService interface {
    Generate(ctx context.Context, reportID string) ([]byte, error)
}
```

El paquete del consumidor importa `addons/reports/port`, no `addons/reports/service` ni `addons/reports/config`. Añade el servicio a su propio `ModuleDeps`, luego su adaptador lo obtiene mediante `modules.RequireService[port.ReportService]`. También agrega el nombre de Reports a `depends` del manifiesto consumidor para que el runtime inicie primero Reports.

No agregues `port/` solo para duplicar interfaces que nadie consume. En Posts, por ejemplo, los contratos compartidos están en `addons/posts/port/ports.go`.

## 5. Generar el registro de fábricas

No edites `internal/modules/backend/factories_gen.go`. El generador busca una función exportada con esta forma en cada `addons/<nombre>/config/`:

```go
func NewReportsModule(cfg ModuleConfig) (*Module, error)
```

Después de crear el manifiesto y el constructor, regenera la lista compilada:

```powershell
go generate ./internal/modules/backend
```

Esto actualiza `internal/modules/backend/factories_gen.go` con imports y fábricas tipadas. El código Go sigue enlazado estáticamente: agrega Reports al ejecutable cuando se genera y compila el registro.

## 6. Crear el adaptador del módulo

Cada módulo conecta su constructor tipado dentro de `config/backend_module.go`. El adaptador recibe la infraestructura común y resuelve contratos publicados por dependencias:

```go
func RegisterBackendModule(registry *modules.Registry) error {
    if err := modules.RegisterFactory(registry, modules.Factory[ModuleConfig, *Module]{
        Name: "reports",
        New:  NewReportsModule,
    }); err != nil {
        return err
    }

    return modules.RegisterBackendInitializer(registry, "reports", func(runtime modules.BackendRuntime) (any, error) {
        auth, err := modules.RequireService[middlewarePort.Authenticator](runtime.Services)
        if err != nil {
            return nil, err
        }

        reportModule, err := NewReportsModule(ModuleConfig{
            AppEnvs:    runtime.AppEnvs,
            AppClients: runtime.AppClients,
            APIv1:      runtime.APIv1,
            DB:         runtime.DB,
            Deps:       ModuleDeps{AuthApiMdw: auth},
        })
        if err != nil {
            return nil, err
        }
        return reportModule, nil
    })
}
```

El orden de llamada no se escribe en código: el backend lo calcula desde `depends`. Si otro módulo consume Reports, publica el contrato después de construirlo:

```go
if err := modules.ProvideService[reportPort.ReportService](runtime.Services, reportModule.ReportService); err != nil {
    return nil, err
}
```

El consumidor declara `"reports"` en su `depends` y pide `reportPort.ReportService` mediante `RequireService`. No edites `internal/modules/backend/composition.go` para conectar módulos. Si `cmd/api/server/v1` necesita directamente un servicio nuevo, expón el módulo en `backend.Modules` y conéctalo en ese límite de la aplicación.

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
| `addons/reports/config/backend_module.go` | Siempre; enlaza el constructor y sus dependencias tipadas con el runtime modular. |
| `internal/modules/backend/factories_gen.go` | Generado con `go generate`; no editar a mano. |
| `internal/modules/backend/composition.go` | No se modifica para agregar módulos. |
| Manifiesto del consumidor | Si otro módulo depende de Reports a nivel de lifecycle. |
| `ModuleDeps` del consumidor y el adaptador del consumidor | Si otro módulo consume un servicio de Reports. |
| `internal/modules/backend/factories_test.go` | Si cambian las expectativas de fábricas oficiales. |
| Pruebas del módulo | Para validar contratos, servicios, persistencia y rutas. |

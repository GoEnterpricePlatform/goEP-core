# Composición de módulos de backend

## Mapa verificado

| Módulo | Consume | Proporciona |
|---|---|---|
| Identity | `AppEnvs`, `AppClients`, MongoDB y multiplexor API v1. Usa los clientes de correo y almacenamiento según configuración. | Servicios de autenticación, usuario, administrador, token, cookie y roles; middleware de autenticación. |
| Posts | `AppEnvs`, `AppClients`, MongoDB, API v1 y el contrato `middlewares/port.Authenticator`. | `posts/port.PostSrv`. |
| Billing | `AppEnvs`, `AppClients`, MongoDB, API v1 y `middlewares/port.Authenticator`. Usa `AppDataPath` para su almacén de claves. | Servicios de proveedores de pago y cifrado desde su constructor; no los consume otro módulo en el bootstrap actual. |
| Catalog | `AppEnvs`, `AppClients`, MongoDB, API v1 y `middlewares/port.Authenticator`. Usa MinIO/Paddle solo cuando están configurados. | Servicios de variaciones y planes desde su constructor; no los consume otro módulo en el bootstrap actual. |
| AI | `AppEnvs`, `AppClients`, API v1 y `posts/port.PostSrv`. Usa el cliente OpenAI cuando está configurado. | `ai-tool-calling/port.ToolCallingSrv`. |
| Web App (fuera del alcance de composición) | Servicios de token, autenticación, cookie, administrador y roles de Identity; Tool Calling de AI; servicio de Posts. | Registra rutas web en los multiplexores existentes. |

`AppClients` agrupa los clientes opcionales de Resend/Gmail, MongoDB, MinIO, OpenAI y Paddle. `InitializeServices` prepara la base de datos y el bucket de MinIO según la configuración. Todos los módulos reciben el contexto tipado que necesitan, aunque algunas configuraciones incluyen campos compartidos no usados por ese módulo.

## Orden y contratos

El orden se calcula desde `depends`: Identity se inicia antes de Posts, Billing y Catalog; Posts se inicia antes de AI. Identity suministra el contrato `middlewares/port.Authenticator`, y Posts suministra `posts/port.PostSrv`. Los adaptadores siguen rellenando `ModuleDeps` tipados dentro de cada paquete. El registro común busca servicios por sus interfaces Go, y los constructores conservan sus tipos concretos.

`internal/modules/backend` consulta el ciclo de vida, recorre el orden calculado desde los manifiestos y devuelve las instancias necesarias al bootstrap. Cada paquete de módulo implementa `RegisterBackendModule` para enlazar su constructor tipado y declarar los contratos que requiere o publica. El adaptador recibe `BackendRuntime`; los contratos entre módulos usan `RequireService[T]` y `ProvideService[T]`. El generador importa los paquetes con manifiesto, por lo que el ejecutable los enlaza estáticamente. Agregar un módulo no requiere editar `composition.go` ni `factories.go`.

El paquete `cmd/api/server/v1` sigue siendo dueño de los recursos globales, `/ping`, `/api/v1/`, middlewares globales y Web App. Las rutas de un módulo habilitado se registran en el mux público; las de un módulo deshabilitado se registran en un mux aislado. Los servicios siguen construyéndose para conservar las dependencias que Web App consume. `backend.Modules` expone explícitamente los módulos cuyos servicios consume el bootstrap web.

## Referencias inspeccionadas

- `cmd/api/server/v1/v1.go`
- `internal/config/app_clients.go`
- `addons/identity/config/module_config.go`
- `addons/posts/config/module_config.go`
- `addons/billing/config/module_config.go`
- `addons/catalog/config/module_config.go`
- `addons/ai-tool-calling/config/module_config.go`
- `addons/identity/config/backend_module.go`
- `addons/posts/config/backend_module.go`
- `addons/billing/config/backend_module.go`
- `addons/catalog/config/backend_module.go`
- `addons/ai-tool-calling/config/backend_module.go`
- Puertos y handlers que consumen autenticación y Posts.

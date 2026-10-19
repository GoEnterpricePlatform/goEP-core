# Hoja de ruta: assets de Web App propiedad de cada módulo

## Decisión propuesta

Sí hay cambios importantes. Los archivos fuente pasarían de estar centralizados en `web/web-app/src/ts` a vivir junto al módulo que los usa, dentro de `addons/<módulo>/static/`. Para hacer funcionar ese cambio hay que actualizar el compilador de TypeScript, la salida que se sirve, la inclusión de archivos en producción y las referencias de las plantillas.

Esta propuesta cubre solo la Web App GoEP del backend. **Next.js queda fuera de este alcance.**

La fase 1 establece el contrato descrito aquí y hace que `cmd/create-module` cree su estructura. Las fases posteriores conectarán el compilador, las URLs de las plantillas y el build de producción.

## Cómo funciona hoy

Actualmente, las fuentes TypeScript están en `web/web-app/src/ts/modules/`. El compilador de [cmd/api/env/dev/ts-compiler/main.go](../cmd/api/env/dev/ts-compiler/main.go) tiene un `entryMap` escrito manualmente y genera JavaScript en `web/web-app/resources/static/js/`.

La Web App sirve `/static/` desde `web/web-app/resources/static/`. En desarrollo usa el directorio del disco; en producción [resources/static_prod.go](../web/web-app/resources/static_prod.go) incluye `static/` en el ejecutable mediante `go:embed`. Las plantillas referencian esas rutas con `resources.StaticPath(...)`.

El compilador TypeScript se ejecuta aparte. El target actual `make run` inicia el servidor Go, pero no ejecuta ese compilador automáticamente.

## Estructura propuesta

Las fuentes y archivos originales pertenecen a cada módulo:

```text
addons/
└── identity/
    ├── manifest.json
    ├── config/
    ├── domain/
    └── static/
        ├── assets.json       # Entradas que se deben compilar
        ├── src/
        │   ├── ts/
        │   └── css/
        └── public/           # Imágenes y archivos que no requieren compilación
```

Los recursos globales siguen en `web/web-app/resources/static/`: favicon, CSS base de la Web App y librerías compartidas. Los archivos generados por el build se escriben bajo una ruta separada por módulo, por ejemplo:

```text
web/web-app/resources/static/
├── assets/favicon.ico
├── css/app.css
├── js/libs/
└── modules/
    └── identity/             # Salida generada; no editar a mano
        ├── js/
        ├── css/
        └── public/
```

Así se centraliza el directorio que el servidor expone y embebe, pero las fuentes editables quedan agrupadas en `addons/<módulo>/static/`. El prefijo `modules/<módulo>/` evita choques entre nombres iguales de archivo.

### Contrato de `static/assets.json`

La fase 1 fija este formato inicial:

```json
{
  "version": 1,
  "scripts": [
    "src/ts/validations/index.ts"
  ],
  "styles": [
    "src/css/identity.css"
  ]
}
```

- `version` identifica el formato del manifiesto de assets.
- `scripts` enumera los puntos de entrada TypeScript, no cada archivo importado por esos puntos de entrada.
- `styles` enumera los puntos de entrada CSS.
- Las rutas son relativas a `addons/<nombre>/static/`, usan `/` y no deben salir de esa carpeta.
- Ambos arreglos pueden estar vacíos, como en el módulo esqueleto.
- Los archivos dentro de `static/public/` son recursos públicos sin compilación y se copian manteniendo la ruta relativa.

La salida acordada conserva el nombre relativo del archivo y agrega una carpeta por módulo. Por ejemplo:

| Fuente | Salida generada | URL pública |
| --- | --- | --- |
| `static/src/ts/validations/index.ts` | `static/modules/identity/js/validations/index.js` | `/static/modules/identity/js/validations/index.js` |
| `static/src/css/identity.css` | `static/modules/identity/css/identity.css` | `/static/modules/identity/css/identity.css` |
| `static/public/images/logo.svg` | `static/modules/identity/public/images/logo.svg` | `/static/modules/identity/public/images/logo.svg` |

En la tabla, `static/` en las rutas de salida representa `web/web-app/resources/static/`. En producción `resources.StaticPath(...)` transformará la ruta a su URL con hash.

## Fases propuestas

### Fase 1: cerrar el contrato de assets por módulo — implementada

- Estructura estándar: `addons/<nombre>/static/{assets.json,src/ts,src/css,public}`.
- Formato `assets.json` versión 1 con listas `scripts` y `styles`.
- Rutas de salida y URLs bajo `modules/<nombre>/`.
- Favicon, CSS base y librerías compartidas permanecen en `web-app/resources/static/`.
- `cmd/create-module` genera las carpetas y un `assets.json` vacío para un módulo nuevo.

**Resultado:** el contrato y el esqueleto ya están definidos. El compilador todavía no descubre ni procesa estos manifiestos.

### Fase 2: descubrir y compilar assets automáticamente — implementada

- Reemplazar el `entryMap` escrito a mano por descubrimiento de archivos TypeScript en el directorio legacy actual.
- Recorrer módulos con `static/`, leer `assets.json` y validar versión, extensiones, ubicación y existencia de las entradas declaradas.
- Compilar TypeScript y CSS declarados a `web/web-app/resources/static/modules/<nombre>/{js,css}/`.
- Copiar `static/public/` a `web/web-app/resources/static/modules/<nombre>/public/`.
- Informar qué módulo o entrada causó un error.
- Mantener temporalmente los TS de `web/web-app/src/ts/` compilándose en sus rutas actuales para no romper las páginas antes de la fase 4.

**Resultado:** agregar un módulo con assets ya no requiere editar una lista global del compilador. Las fuentes de Identity aún no se han movido; eso ocurrirá en la fase 4.

### Fase 3: conectar desarrollo y producción

- Asegurar que el servidor de desarrollo lea la salida generada desde `resources/static/`.
- Ejecutar el compilador antes de iniciar la Web App, o documentar un comando/watch de desarrollo que lo mantenga actualizado.
- Vigilar cambios en `static/public/` durante el modo `-watch` (por ahora se copian al iniciar el compilador).
- Asegurar que la build de producción genere los assets antes de compilar Go, porque `go:embed` solo incluye los archivos presentes durante la compilación.
- Mantener las URLs hash de producción que genera `hashfs` mediante `resources.StaticPath(...)`.
- Limpiar la salida generada antigua antes de regenerar para que no se sirvan archivos obsoletos de módulos eliminados.

**Resultado:** los mismos assets y URLs funcionan en desarrollo y en el ejecutable de producción.

### Fase 4: actualizar plantillas y ciclo de vida

- Mover los TS actuales de `web/web-app/src/ts/modules/identity/` a `addons/identity/static/src/` y declarar sus entradas.
- Actualizar las plantillas Templ para usar las nuevas rutas por módulo mediante `resources.StaticPath(...)`.
- Determinar cómo se comportan los assets de un módulo deshabilitado. Los assets públicos no deben tratarse como secretos; si el requisito es que no se sirvan mientras el módulo está deshabilitado, el servidor estático debe consultar el estado del módulo.
- Actualizar y probar el flujo de creación de un módulo con y sin assets.

**Resultado:** Identity usa la nueva estructura y las futuras plantillas de módulos siguen el mismo contrato.

## Qué significa instalar o desinstalar

Poner fuentes en `addons/<módulo>/static/` organiza y entrega propiedad de esos archivos al módulo. No hace que Go cargue o elimine código en caliente.

- Agregar o quitar un módulo del código requiere regenerar el registro backend y reconstruir la aplicación.
- Los assets generados quedan incluidos en la build de producción.
- El ciclo de vida de la aplicación puede habilitar o deshabilitar módulos ya compilados. Si se requiere que un módulo deshabilitado tampoco sirva sus archivos estáticos, eso debe implementarse como una regla del servidor estático.
- No se guardan los archivos de assets en MongoDB.

## Cambios involucrados

La migración afecta varios puntos del repositorio:

| Área | Cambio |
| --- | --- |
| `cmd/api/env/dev/ts-compiler` | Descubrimiento de manifiestos de assets y compilación por módulo |
| `cmd/create-module` | Esqueleto estático para módulos nuevos |
| `web/web-app/resources` | Mantener el servicio común y el embedding de la salida generada |
| Plantillas Templ | Actualizar las URLs de scripts, estilos e imágenes de módulo |
| `Makefile` y flujo de build | Ejecutar compilación de assets antes de servir o compilar el binario |
| Ciclo de vida | Decidir si los assets de módulos deshabilitados siguen siendo accesibles |

Por eso conviene implementarlo por fases y migrar un módulo primero. La primera migración propuesta es Identity, porque ya tiene páginas y TypeScript propios.

## Fuera de alcance

- Assets de Next.js.
- Cambiar el favicon o eliminar los recursos compartidos de Web App.
- Guardar assets en MongoDB.
- Cargar paquetes Go dinámicamente sin reconstruir el ejecutable.

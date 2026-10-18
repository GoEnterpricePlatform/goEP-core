# Persistencia del ciclo de vida de módulos

El estado de instalación corresponde a toda la instalación de GoEP Core. No existe un modelo de organización/tenant ni un identificador de tenant en las rutas o repositorios revisados; los módulos comparten una única base de datos Mongo configurada por la aplicación.

El código del ciclo de vida está en `internal/modules/lifecycle/`, separado en estado y contrato de almacenamiento (`state.go`), operaciones del gestor (`manager.go`) y persistencia Mongo (`mongo_store.go`). El adaptador Mongo guarda un documento por módulo en `goep_core_modules`, identificado por un índice único en `module`. Los manifiestos determinan disponibilidad y versión; la ausencia de un documento significa que el módulo está disponible pero no instalado. Los campos `installed` y `enabled` distinguen los otros estados.

`Manager.Install` crea o actualiza el estado sin activar el módulo. Repetir instalación no altera un estado de igual versión; al cambiar la versión del manifiesto actualiza la versión persistida y conserva la activación existente. `Manager.Enable` requiere instalación y dependencias directas y transitivas instaladas y habilitadas. `Manager.Disable` desactiva el módulo sin borrar su estado ni tocar las colecciones de datos, y rechaza desactivar un módulo requerido por otro módulo habilitado.

La desinstalación no tiene operación en esta fase. La política definida es permitirla solo después de deshabilitar el módulo y eliminar únicamente su documento de estado; nunca borrar automáticamente las colecciones de datos del módulo.

El adaptador crea el índice único mediante `lifecycle.MongoStore.CreateIndexes(ctx)`; la futura inicialización del Module Manager debe llamarlo antes de operar. El adaptador no se conecta aún al bootstrap HTTP; la integración de lectura de estado y activación por rutas corresponde a la fase de integración del Module Manager.

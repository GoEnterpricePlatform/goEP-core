# Module lifecycle persistence

The installation status corresponds to the entire GoEP Core installation. There is no tenant/organization model or tenant identifier in the routes or repositories reviewed; The modules share a single Mongo database configured by the application.

The lifecycle code is in `internal/modules/lifecycle/`, separated into state and storage contract (`state.go`), manager operations (`manager.go`) and Mongo persistence (`mongo_store.go`). The Mongo adapter stores one document per module in `goep_core_modules`, identified by a unique index in `module`. Manifests determine availability and version; The absence of a document means that the module is available but not installed. The `installed` and `enabled` fields distinguish the other states.

`Manager.Install` creates or updates the state without activating the module. Repeating installation does not alter a state of the same version; changing the manifest version updates the persisted version and preserves the existing activation. `Manager.Enable` requires installation and direct and transitive dependencies installed and enabled. `Manager.Disable` disables the module without clearing its state or touching data collections, and refuses to disable a module required by another enabled module.

The uninstall has no operation in this phase. The defined policy is to allow it only after disabling the module and deleting only its state document; never automatically delete module data collections.

The backend bootstrap creates the unique index before querying the states. In an empty base, install and enable the included modules to preserve the initial behavior. On subsequent boots it uses the persisted state to expose or isolate its API routes. As long as Web App is still out of the lifecycle, bootstrap builds backend services even when their API routes are disabled, because Web App still consumes them directly.
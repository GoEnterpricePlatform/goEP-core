# Log compiled deployments

`internal/modules.Registry` maps the names declared in the manifests to the factories that are compiled into the executable. Each factory maintains specific types of configuration and output for its module; the registry does not expose a common service container.

`internal/modules/backend/factories_gen.go` is generated from the manifests under `addons/` and the conventional typed constructors in each module's `config` package. Run `go generate ./internal/modules/backend` after adding a module. The generator fails if a module does not have exactly one exported constructor shaped as `func(ModuleConfig) (*Module, error)`.

`Registry.Validate` compares discovered manifests with the registry in both directions: reporting a discovered module with no compiled implementation and a compiled factory with no manifest. The backend bootstrap runs discovery, registry validation and dependency resolution before constructing modules.

`internal/modules.NewDependencyResolver` validates that declared dependencies exist, checks compiled availability and detects cycles. `Order` returns manifests in deterministic topological order. `ValidateCanEnable` checks direct and transitive dependencies against the modules already enabled. `internal/modules/backend.Initialize` walks that order and calls each module's local `RegisterBackendModule` adapter. The adapter builds its own typed configuration and uses `RequireService[T]` and `ProvideService[T]` for cross-module contracts. New modules no longer require a change to the central composition file.

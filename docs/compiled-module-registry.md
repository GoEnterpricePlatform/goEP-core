# Log compiled deployments

`internal/modules.Registry` maps the names declared in the manifests to the factories that are compiled into the executable. Each factory maintains specific types of configuration and output for its module; The registry does not expose a common service container.

`internal/modules/backend.NewCompiledRegistry` registers Identity, Posts, Billing, Catalog and AI Tool Calling. The composer uses those factories along with the typed injection defined in Phase 1. Adding an official Go implementation requires registering your factory and generating a new build.

`Registry.Validate` compares discovered manifests with the registry in both directions: reporting a discovered module with no compiled implementation and a compiled factory with no manifest. This validation is available for subsequent Module Manager integration; discovery is not connecting to HTTP boot yet.
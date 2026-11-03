# Module discovery

`internal/modules.Discovery` reads the immediate subdirectories of an addons folder and validates each `manifest.json`. `Discover()` returns backend modules sorted by name; `DiscoverAll()` returns both backend and web modules. Manifests without `type` or with `"type": "backend"` are backend modules. Modules with `"type": "web"` are initialized by the web composition stage after backend services are available, and are excluded from the backend lifecycle.

The policy for a subdirectory without `manifest.json` is to fail discovery with an error indicating incomplete folder. Invalid manifests and duplicate module names are also reported. Files other than directories are ignored.

The discoverer only reads metadata: it does not load or execute `.go` files, it does not record paths, and it does not determine whether a module's code is compiled into the executable.

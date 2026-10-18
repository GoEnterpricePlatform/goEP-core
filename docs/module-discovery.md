# Module discovery

`internal/modules.Discovery` reads the immediate subdirectories of an addons folder, parses each `manifest.json` with the manifest validator, and returns the results sorted by module name.

The policy for a subdirectory without `manifest.json` is to fail discovery with an error indicating incomplete folder. Invalid manifests and duplicate module names are also reported. Files other than directories are ignored.

The discoverer only reads metadata: it does not load or execute `.go` files, it does not record paths, and it does not determine whether a module's code is compiled into the executable.

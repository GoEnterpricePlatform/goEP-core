package backend

// Regenerate module initializers and per-addon static asset packages after changing addons.
//go:generate go run -buildvcs=false ../../../cmd/goep generate-module-registry -addons ../../../addons -out factories_gen.go -assets-out ../assets/static_assets_gen.go

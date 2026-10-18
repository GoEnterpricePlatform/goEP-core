package backend

import (
	aiModule "github.com/GoEnterpricePlatform/goEP-core/addons/ai-tool-calling/config"
	billingModule "github.com/GoEnterpricePlatform/goEP-core/addons/billing/config"
	catalogModule "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/config"
	identityModule "github.com/GoEnterpricePlatform/goEP-core/addons/identity/config"
	postModule "github.com/GoEnterpricePlatform/goEP-core/addons/posts/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
)

// NewCompiledRegistry registers the backend module constructors linked into
// this executable. Adding a Go module still requires registering its factory
// here and rebuilding the binary.
func NewCompiledRegistry() (*modules.Registry, error) {
	registry := modules.NewRegistry()

	if err := modules.RegisterFactory(registry, modules.Factory[identityModule.ModuleConfig, *identityModule.Module]{
		Name: "identity",
		New:  identityModule.NewIdentityModule,
	}); err != nil {
		return nil, err
	}
	if err := modules.RegisterFactory(registry, modules.Factory[postModule.ModuleConfig, *postModule.Module]{
		Name: "posts",
		New:  postModule.NewPostModule,
	}); err != nil {
		return nil, err
	}
	if err := modules.RegisterFactory(registry, modules.Factory[billingModule.ModuleConfig, *billingModule.Module]{
		Name: "billing",
		New:  billingModule.NewBillingModule,
	}); err != nil {
		return nil, err
	}
	if err := modules.RegisterFactory(registry, modules.Factory[catalogModule.ModuleConfig, *catalogModule.Module]{
		Name: "catalog",
		New:  catalogModule.NewCatalogModule,
	}); err != nil {
		return nil, err
	}
	if err := modules.RegisterFactory(registry, modules.Factory[aiModule.ModuleConfig, *aiModule.Module]{
		Name: "ai-tool-calling",
		New:  aiModule.NewAiModule,
	}); err != nil {
		return nil, err
	}

	return registry, nil
}

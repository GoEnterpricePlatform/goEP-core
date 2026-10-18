// Package backend owns the composition of the statically linked backend modules.
// It keeps constructor and dependency wiring details out of the HTTP bootstrap.
package backend

import (
	"net/http"

	aiModule "github.com/GoEnterpricePlatform/goEP-core/addons/ai-tool-calling/config"
	billingModule "github.com/GoEnterpricePlatform/goEP-core/addons/billing/config"
	catalogModule "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/config"
	identityModule "github.com/GoEnterpricePlatform/goEP-core/addons/identity/config"
	postModule "github.com/GoEnterpricePlatform/goEP-core/addons/posts/config"
	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/lifecycle"
)

// Config contains the infrastructure shared by backend modules and their route
// registrar. Each module still receives its own typed configuration.
type Config struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	AddonsDir  string
}

// Modules contains the backend modules and services needed by the web layer.
type Modules struct {
	Identity  *identityModule.Module
	Posts     *postModule.Module
	AI        *aiModule.Module
	Lifecycle *lifecycle.Manager
}

// Initialize constructs backend modules in dependency order and registers
// their API routes on cfg.APIv1.
func Initialize(cfg Config) (*Modules, error) {
	registry, enabled, lifecycleManager, err := initializeModuleLifecycle(cfg)
	if err != nil {
		return nil, err
	}
	identityFactory, err := modules.GetFactory[identityModule.ModuleConfig, *identityModule.Module](registry, "identity")
	if err != nil {
		return nil, err
	}
	identity, err := identityFactory.New(identityModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      moduleRoutes(cfg.APIv1, enabled, "identity"),
		DB:         cfg.AppClients.DB,
	})
	if err != nil {
		return nil, err
	}

	postsFactory, err := modules.GetFactory[postModule.ModuleConfig, *postModule.Module](registry, "posts")
	if err != nil {
		return nil, err
	}
	posts, err := postsFactory.New(postModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      moduleRoutes(cfg.APIv1, enabled, "posts"),
		DB:         cfg.AppClients.DB,
		Deps: postModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	billingFactory, err := modules.GetFactory[billingModule.ModuleConfig, *billingModule.Module](registry, "billing")
	if err != nil {
		return nil, err
	}
	_, err = billingFactory.New(billingModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      moduleRoutes(cfg.APIv1, enabled, "billing"),
		DB:         cfg.AppClients.DB,
		Deps: billingModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	catalogFactory, err := modules.GetFactory[catalogModule.ModuleConfig, *catalogModule.Module](registry, "catalog")
	if err != nil {
		return nil, err
	}
	_, err = catalogFactory.New(catalogModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      moduleRoutes(cfg.APIv1, enabled, "catalog"),
		DB:         cfg.AppClients.DB,
		Deps: catalogModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	aiFactory, err := modules.GetFactory[aiModule.ModuleConfig, *aiModule.Module](registry, "ai-tool-calling")
	if err != nil {
		return nil, err
	}
	ai, err := aiFactory.New(aiModule.ModuleConfig{
		AppEnvs:     cfg.AppEnvs,
		AppClients:  cfg.AppClients,
		APIv1:       moduleRoutes(cfg.APIv1, enabled, "ai-tool-calling"),
		PostService: posts.PostService,
	})
	if err != nil {
		return nil, err
	}

	return &Modules{
		Identity:  identity,
		Posts:     posts,
		AI:        ai,
		Lifecycle: lifecycleManager,
	}, nil
}

func moduleRoutes(apiV1 *http.ServeMux, enabled map[string]bool, moduleName string) *http.ServeMux {
	if enabled[moduleName] {
		return apiV1
	}
	// Backend services are still constructed during this transition because
	// Web App consumes them directly. A private mux keeps disabled API routes
	// out of the public router until Web App joins the lifecycle manager.
	return http.NewServeMux()
}

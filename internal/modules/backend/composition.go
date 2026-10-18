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
)

// Config contains the infrastructure shared by backend modules and their route
// registrar. Each module still receives its own typed configuration.
type Config struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
}

// Modules contains the backend modules and services needed by the web layer.
type Modules struct {
	Identity *identityModule.Module
	Posts    *postModule.Module
	AI       *aiModule.Module
}

// Initialize constructs backend modules in dependency order and registers
// their API routes on cfg.APIv1.
func Initialize(cfg Config) (*Modules, error) {
	identity, err := identityModule.NewIdentityModule(identityModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      cfg.APIv1,
		DB:         cfg.AppClients.DB,
	})
	if err != nil {
		return nil, err
	}

	posts, err := postModule.NewPostModule(postModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      cfg.APIv1,
		DB:         cfg.AppClients.DB,
		Deps: postModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = billingModule.NewBillingModule(billingModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      cfg.APIv1,
		DB:         cfg.AppClients.DB,
		Deps: billingModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = catalogModule.NewCatalogModule(catalogModule.ModuleConfig{
		AppEnvs:    cfg.AppEnvs,
		AppClients: cfg.AppClients,
		APIv1:      cfg.APIv1,
		DB:         cfg.AppClients.DB,
		Deps: catalogModule.ModuleDeps{
			AuthApiMdw: identity.AuthApiMdw,
		},
	})
	if err != nil {
		return nil, err
	}

	ai, err := aiModule.NewAiModule(aiModule.ModuleConfig{
		AppEnvs:     cfg.AppEnvs,
		AppClients:  cfg.AppClients,
		APIv1:       cfg.APIv1,
		PostService: posts.PostService,
	})
	if err != nil {
		return nil, err
	}

	return &Modules{
		Identity: identity,
		Posts:    posts,
		AI:       ai,
	}, nil
}

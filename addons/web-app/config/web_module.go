package config

import (
	tcPort "github.com/GoEnterpricePlatform/goEP-core/addons/ai-tool-calling/port"
	adminPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/port"
	authPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	rolePort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	tokenPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	postPort "github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	cookiePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/handler/cookie/port"
)

// RegisterWebModule connects Web App to the modular web bootstrap.
func RegisterWebModule(registry *modules.Registry) error {
	return modules.RegisterWebInitializer(registry, "web-app", initializeWebApp)
}

func initializeWebApp(runtime modules.WebRuntime) error {
	tokenSrv, err := modules.RequireService[tokenPort.TokenSrv](runtime.Services)
	if err != nil {
		return err
	}
	authSrv, err := modules.RequireService[authPort.AuthSrv](runtime.Services)
	if err != nil {
		return err
	}
	cookieSrv, err := modules.RequireService[cookiePort.CookieSrv](runtime.Services)
	if err != nil {
		return err
	}
	adminSrv, err := modules.RequireService[adminPort.AdminSrv](runtime.Services)
	if err != nil {
		return err
	}
	tcService, _ := modules.RequireService[tcPort.ToolCallingSrv](runtime.Services)
	postService, _ := modules.RequireService[postPort.PostSrv](runtime.Services)
	rolesSrv, err := modules.RequireService[rolePort.RoleSrv](runtime.Services)
	if err != nil {
		return err
	}

	NewWebAppModule(ModuleConfig{
		AppEnvs:    runtime.AppEnvs,
		Mux:        runtime.Mux,
		TemplateV1: runtime.TemplateV1,
		Deps: ModuleDeps{
			TokenSrv:    tokenSrv,
			AuthSrv:     authSrv,
			CookieSrv:   cookieSrv,
			AdminSrv:    adminSrv,
			TCService:   tcService,
			PostService: postService,
			RolesSrv:    rolesSrv,
		},
		AvailableModules: runtime.AvailableModules,
	})
	return nil
}

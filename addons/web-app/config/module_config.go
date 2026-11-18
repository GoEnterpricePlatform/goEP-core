package config

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"

	adminP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/port"
	authP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/handler/cookie/port"

	roleP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	identityUIHandler "github.com/GoEnterpricePlatform/goEP-core/addons/identity/ui/handler"
	adminHandler "github.com/GoEnterpricePlatform/goEP-core/addons/web-app/modules/identity/admin/api/handler"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/resources"
	srcHandler "github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/api/handler"
	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/layouts"
	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

type ModuleConfig struct {
	AppEnvs          *config.AppEnvs
	Mux              *http.ServeMux
	TemplateV1       *http.ServeMux
	Deps             ModuleDeps
	AvailableModules []modules.Manifest
}

// ModuleDeps groups the external dependencies (from other modules/services)
// that this module needs to function.
type ModuleDeps struct {
	TokenSrv  tokenP.TokenSrv
	AuthSrv   authP.AuthSrv
	CookieSrv cookieP.CookieSrv
	AdminSrv  adminP.AdminSrv
	RolesSrv  roleP.RoleSrv
}

func NewWebAppModule(cfg ModuleConfig) {

	cfg.Mux.Handle("/static/", resources.Handler())

	mdwSrvTmpl := middlewareServiceTmpl.NewMdwSrvTmpl(cfg.Deps.TokenSrv, cfg.Deps.AuthSrv, cfg.Deps.CookieSrv)

	// Templates - admin
	adminH := adminHandler.NewAdminHandler(cfg.Deps.AdminSrv, cfg.Deps.CookieSrv, cfg.AppEnvs.ApiBaseUrl, mdwSrvTmpl)
	adminH.RegisterRoutes(cfg.Mux, cfg.TemplateV1)

	srcHandler.NewSrcHandler(cfg.Mux, cfg.TemplateV1, mdwSrvTmpl, cfg.AvailableModules)

	identityUI := identityUIHandler.NewHandler(
		cfg.Deps.RolesSrv,
		cfg.AppEnvs.ApiBaseUrl,
		mdwSrvTmpl,
		layouts.DashboardLayout,
	)
	identityUI.RegisterRoutes(cfg.TemplateV1)

}

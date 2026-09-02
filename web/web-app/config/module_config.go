package config

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"

	adminP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/port"
	authP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/port"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/port"

	tcP "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"
	postP "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"

	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
	adminHandler "github.com/GoEnterpricePlatform/goEP-core/web/web-app/modules/identity/admin/api/handler"

	/* postHandlerWeb "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/admin/api/posts/handler"
	adminRenderer "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/admin/renderer"
	chatToolCallingHandler "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/ai-tool-calling/api/handler"
	toolCallingHandler "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/ai-tool-calling/api/posts/handler"
	toolCallingRenderer "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/ai-tool-calling/renderer" */
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/resources"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/src/api/handler"
)

type ModuleConfig struct {
	AppEnvs    *config.AppEnvs
	Mux        *http.ServeMux
	TemplateV1 *http.ServeMux
	Deps       ModuleDeps
}

// ModuleDeps groups the external dependencies (from other modules/services)
// that this module needs to function.
type ModuleDeps struct {
	TokenSrv    tokenP.TokenSrv
	AuthSrv     authP.AuthSrv
	CookieSrv   cookieP.CookieSrv
	AdminSrv    adminP.AdminSrv
	TCService   tcP.ToolCallingSrv
	PostService postP.PostSrv
}

func NewWebAppModule(cfg ModuleConfig) {

	cfg.Mux.Handle("/static/", resources.Handler())

	mdwSrvTmpl := middlewareServiceTmpl.NewMdwSrvTmpl(cfg.Deps.TokenSrv, cfg.Deps.AuthSrv, cfg.Deps.CookieSrv)

	// Templates - admin
	adminH := adminHandler.NewAdminHandler(cfg.Deps.AdminSrv, cfg.Deps.CookieSrv, cfg.AppEnvs.ApiBaseUrl, mdwSrvTmpl)
	adminH.RegisterRoutes(cfg.Mux, cfg.TemplateV1)

	handler.NewSrcHandler(cfg.Mux, cfg.TemplateV1, cfg.Deps.PostService, mdwSrvTmpl)

	// Templates - admin
	// adminH := adminHandler.NewAdminHandler(cfg.Deps.AdminSrv, cfg.Deps.CookieSrv, cfg.AppEnvs.ApiBaseUrl, adminR, mdwSrvTmpl)
	// adminH.RegisterRoutes(cfg.Mux, cfg.TemplateV1)

	// Templates - chat tool calling
	// toolCalllingR := toolCallingRenderer.NewToolCallingRenderer()
	// toolCallingH := chatToolCallingHandler.NewChatToolCallingHandler(
	// 	cfg.Deps.TCService,
	// 	cfg.AppEnvs.ApiBaseUrl,
	// 	toolCalllingR,
	// 	mdwSrvTmpl,
	// 	cfg.Deps.PostService,
	// )
	// toolCallingH.RegisterRoutes(cfg.Mux, cfg.TemplateV1)

	// tool Calling posts
	// toolCallingPostsH := toolCallingHandler.NewPostWebAiHandler(
	// 	cfg.Deps.PostService,
	// 	cfg.AppEnvs.ApiBaseUrl,
	// 	toolCalllingR,
	// 	mdwSrvTmpl,
	// )
	// toolCallingPostsH.RegisterRoutes(cfg.TemplateV1)

	// Templates - post
	// postH := postHandlerWeb.NewPostWebHandler(cfg.Deps.PostService, adminH.ApiBaseUrl, adminR)
	// postH.RegisterRoutes(cfg.TemplateV1)

	// Templates - public
	//publicHandler.NewPublicHandler(cfg.Mux, cfg.AppEnvs.ApiBaseUrl, cfg.Deps.PostService)

}

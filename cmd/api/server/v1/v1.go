package v1

import (
	"log"
	"net/http"

	tcPort "github.com/GoEnterpricePlatform/goEP-core/addons/ai-tool-calling/port"
	adminPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/port"
	authPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	rolePort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	tokenPort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	postPort "github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	webAppModule "github.com/GoEnterpricePlatform/goEP-core/web/web-app/config"

	backendModules "github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend"

	cookiePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/handler/cookie/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/logger"
	sharedHandler "github.com/GoEnterpricePlatform/goEP-core/internal/shared/handler"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// get envs
	appEnvs := appConfig.NewAppEnvs()
	appEnvs.Load()

	appClients := appConfig.NewClients()
	appClients.GetClients(appEnvs)
	appClients.InitializeServices(appEnvs)

	zapLogger := logger.NewHttpLogger(appEnvs.AppEnv)

	corsLogger := middlewares.CorsMiddleware(appEnvs.AllowedOrigins)

	// Add global middlewares, the order matters (logger → CORS → router)
	apiHandler := zapLogger.Middleware(corsLogger(mux))

	// Api version
	// Note: all subsequent handlers should also be registered using v1
	v1 := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// It is used to know if the server is running
	mux.HandleFunc("GET /ping", sharedHandler.Ping)

	// Call the modules

	backend, err := backendModules.Initialize(backendModules.Config{
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		AddonsDir:  "addons",
	})
	if err != nil {
		log.Fatal(err)
	}

	tokenSrv, err := modules.RequireService[tokenPort.TokenSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	authSrv, err := modules.RequireService[authPort.AuthSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	cookieSrv, err := modules.RequireService[cookiePort.CookieSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	adminSrv, err := modules.RequireService[adminPort.AdminSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	tcService, err := modules.RequireService[tcPort.ToolCallingSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	postService, err := modules.RequireService[postPort.PostSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}
	rolesSrv, err := modules.RequireService[rolePort.RoleSrv](backend.Services)
	if err != nil {
		log.Fatal(err)
	}

	// web pages do not have the /api prefix
	templateV1 := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", templateV1))

	// Web modules

	// for import module :
	// standardWebModule "github.com/GoEnterpricePlatform/goEP-core/web/standard-web/config"
	/* standardWebModule.NewStandardWebModule(standardWebModule.ModuleConfig{
		AppEnvs: appEnvs,
		Mux: mux,
		TemplateV1: templateV1,
		Deps: standardWebModule.ModuleDeps{
			TokenSrv: identityMdl.TokenSrv,
			AuthSrv: identityMdl.AuthSrv,
			CookieSrv: identityMdl.CookieSrv,
			AdminSrv: identityMdl.AdminSrv,
			TCService: aiMdl.TCService,
			PostService: postMdl.PostService,
		},
	}) */

	webAppModule.NewWebAppModule(webAppModule.ModuleConfig{
		AppEnvs:    appEnvs,
		Mux:        mux,
		TemplateV1: templateV1,
		Deps: webAppModule.ModuleDeps{
			TokenSrv:    tokenSrv,
			AuthSrv:     authSrv,
			CookieSrv:   cookieSrv,
			AdminSrv:    adminSrv,
			TCService:   tcService,
			PostService: postService,
			RolesSrv:    rolesSrv,
		},
		AvailableModules: backend.AvailableModules,
	})

	return apiHandler
}

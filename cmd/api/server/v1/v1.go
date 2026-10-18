package v1

import (
	"log"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	webAppModule "github.com/GoEnterpricePlatform/goEP-core/web/web-app/config"

	backendModules "github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend"

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
	})
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
			TokenSrv:    backend.Identity.TokenSrv,
			AuthSrv:     backend.Identity.AuthSrv,
			CookieSrv:   backend.Identity.CookieSrv,
			AdminSrv:    backend.Identity.AdminSrv,
			TCService:   backend.AI.TCService,
			PostService: backend.Posts.PostService,
			RolesSrv:    backend.Identity.RoleSrv,
		},
	})

	return apiHandler
}

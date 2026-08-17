package v1

import (
	"log"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	webAppModule "github.com/GoEnterpricePlatform/goEP-core/web/web-app/config"

	aiModule "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/config"
	billingModule "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/config"
	catalogModule "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/config"

	identityModule "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/config"
	postModule "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/config"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares/logger"
	sharedHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/handler"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// define app Stack
	appStack := appConfig.NewAppStack(appConfig.DBMongo, appConfig.MailGmail, appConfig.FSMinio)

	// get envs
	appEnvs := appConfig.NewAppEnvs()
	appEnvs.Load(appStack)

	appClients := appConfig.NewClients()
	appClients.GetClients(appStack, appEnvs)

	// Initialize the infrastructure in the services that we are using and that are strictly
	// necessary for all modules; if it is only for a specific module, it should go in the
	// initializer folder within your module.
	// TODO: create an Initialize() method for AppClients. Both services are required so at this time
	// TODO:there is no validation they both run.
	// - MongoDB
	db := appClients.MongoConn.DB.Database(appEnvs.MongoInitDB)
	appClients.MongoConn.Ping()

	// - Minio
	err := appClients.MinioCli.CreateStorage(appEnvs.MinioBucketName)
	if err != nil {
		log.Fatal(err)
	}

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

	identityMdl, err := identityModule.NewIdentityModule(identityModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		DB:         db,
	})
	if err != nil {
		log.Fatal(err)
	}

	postMdl, err := postModule.NewPostModule(postModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		DB:         db,
		Deps: postModule.ModuleDeps{
			AuthApiMdw: identityMdl.AuthApiMdw,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = billingModule.NewBillingModule(billingModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		DB:         db,
		Deps: billingModule.ModuleDeps{
			AuthApiMdw: identityMdl.AuthApiMdw,
		},
	})

	_, err = catalogModule.NewCatalogModule(catalogModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		DB:         db,
		Deps: catalogModule.ModuleDeps{
			AuthApiMdw: identityMdl.AuthApiMdw,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Optional services
	aiMdl, err := aiModule.NewAiModule(aiModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		PostModule: postMdl,
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
			TokenSrv:    identityMdl.TokenSrv,
			AuthSrv:     identityMdl.AuthSrv,
			CookieSrv:   identityMdl.CookieSrv,
			AdminSrv:    identityMdl.AdminSrv,
			TCService:   aiMdl.TCService,
			PostService: postMdl.PostService,
		},
	})

	return apiHandler
}

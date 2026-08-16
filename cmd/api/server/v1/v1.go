package v1

import (
	"log"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"

	aiModule "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/config"
	billingModule "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/config"
	catalogModule "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/config"

	identityModule "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/config"
	postModule "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/config"
	adminHandler "github.com/GoEnterpricePlatform/goEP-core/web/admin/api/handler"
	postHandlerWeb "github.com/GoEnterpricePlatform/goEP-core/web/admin/api/posts/handler"
	adminRenderer "github.com/GoEnterpricePlatform/goEP-core/web/admin/renderer"
	chatToolCallingHandler "github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/api/handler"
	toolCallingHandler "github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/api/posts/handler"
	toolCallingRenderer "github.com/GoEnterpricePlatform/goEP-core/web/ai-tool-calling/renderer"
	publicHandler "github.com/GoEnterpricePlatform/goEP-core/web/public/api/handler"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares/logger"
	sharedHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/handler"

	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
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
	})

	_, err = catalogModule.NewCatalogModule(catalogModule.ModuleConfig{
		AppStack:   appStack,
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		DB:         db,
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

	// Templates
	adminR := adminRenderer.NewAdminRenderer()
	mdwSrvTmpl := middlewareServiceTmpl.NewMdwSrvTmpl(identityMdl.TokenSrv, identityMdl.AuthSrv, identityMdl.CookieSrv)

	// no lleva prefix api
	templateV1 := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", templateV1))

	// Templates - admin
	adminH := adminHandler.NewAdminHandler(identityMdl.AdminSrv, identityMdl.CookieSrv, appEnvs.ApiBaseUrl, adminR, mdwSrvTmpl)
	adminH.RegisterRoutes(mux, templateV1)

	// Templates - chat tool calling
	toolCalllingR := toolCallingRenderer.NewToolCallingRenderer()
	toolCallingH := chatToolCallingHandler.NewChatToolCallingHandler(aiMdl.TCService, appEnvs.ApiBaseUrl, toolCalllingR, mdwSrvTmpl, postMdl.PostService)
	toolCallingH.RegisterRoutes(mux, templateV1)

	// tool Calling posts

	toolCallingPostsH := toolCallingHandler.NewPostWebAiHandler(postMdl.PostService, appEnvs.ApiBaseUrl, toolCalllingR, mdwSrvTmpl)
	toolCallingPostsH.RegisterRoutes(templateV1)

	// Templates - post
	postH := postHandlerWeb.NewPostWebHandler(postMdl.PostService, adminH.ApiBaseUrl, adminR)
	postH.RegisterRoutes(templateV1)

	// Templates - public
	publicHandler.NewPublicHandler(mux, appEnvs.ApiBaseUrl, postMdl.PostService)

	return apiHandler
}

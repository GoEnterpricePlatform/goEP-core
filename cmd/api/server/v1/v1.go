package v1

import (
	"log"
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	backendModules "github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend"

	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/logger"
	sharedHandler "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/handler"
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

	// Web pages do not have the /api prefix.
	templateV1 := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", templateV1))

	_, err := backendModules.Initialize(backendModules.Config{
		AppEnvs:    appEnvs,
		AppClients: appClients,
		APIv1:      v1,
		WebMux:     mux,
		TemplateV1: templateV1,
		AddonsDir:  "addons",
	})
	if err != nil {
		log.Fatal(err)
	}

	return apiHandler
}

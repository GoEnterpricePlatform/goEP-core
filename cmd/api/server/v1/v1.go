package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/amorindev/go-tmpl/internal/config"
	minioClient "github.com/amorindev/go-tmpl/internal/minio"
	mongoClient "github.com/amorindev/go-tmpl/internal/mongo"
	"github.com/amorindev/go-tmpl/internal/tokens/service"
	adminHandler "github.com/amorindev/go-tmpl/pkg/app/admin/api/handler"
	authMethodHandler "github.com/amorindev/go-tmpl/pkg/app/auth-methods/handler"
	authMethodService "github.com/amorindev/go-tmpl/pkg/app/auth-methods/service"
	roleInitializer "github.com/amorindev/go-tmpl/pkg/app/roles/initializer"
	roleRepository "github.com/amorindev/go-tmpl/pkg/app/roles/repository/mongo"
	sessionRepository "github.com/amorindev/go-tmpl/pkg/app/session/repository/mongo"
	sessionService "github.com/amorindev/go-tmpl/pkg/app/session/service"
	userFileStorage "github.com/amorindev/go-tmpl/pkg/app/users/file-storage/minio"
	userHandler "github.com/amorindev/go-tmpl/pkg/app/users/handler"
	userRepository "github.com/amorindev/go-tmpl/pkg/app/users/repository/mongo"
	userService "github.com/amorindev/go-tmpl/pkg/app/users/service"
	"github.com/amorindev/go-tmpl/pkg/shared/api/middlewares"
)

func New() http.Handler {
	mux := http.NewServeMux()

	appEnvs := config.Load()

	// Add global middlewares
	muxWithCors := middlewares.CorsMiddleware(appEnvs.AllowedOrigins)(mux)

	// Api version
	v1 := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	// MongoDB
	mongoConn := mongoClient.New(appEnvs.MongoDBUri)
	mongoDB := mongoConn.DB.Database(appEnvs.MongoInitDB)
	mongoConn.Ping()

	// Minio
	minioC, err := minioClient.NewClient(appEnvs.MinioEndpoint, appEnvs.MinioAccessKey, appEnvs.MinioSecretKey, appEnvs.MinioUseSSL)
	if err != nil {
		log.Fatal(err)
	}

	err = minioC.CreateStorage(appEnvs.MinioBucketName)
	if err != nil {
		log.Fatal(err)
	}

	// Collections
	userColl := mongoDB.Collection("users")
	sessionColl := mongoDB.Collection("sessions")
	roleColl := mongoDB.Collection("roles")

	// Repositories
	userRepo := userRepository.NewUserRepo(mongoConn.DB, userColl)
	sessionRepo := sessionRepository.NewSessionRepo(mongoConn.DB, sessionColl)
	roleRepo := roleRepository.NewRoleRepo(mongoConn.DB, roleColl)

	// Indexes
	err = userRepo.CreateIndexes()
	if err != nil {
		log.Fatal(err)
	}

	// File Storage
	userFileStg := userFileStorage.NewUserFileStg(minioC.Client, appEnvs.MinioBucketName, 0)

	// Services
	authSrv := service.NewTokenSrv(appEnvs.JWTAccessSecret, appEnvs.JWTRefreshSecret, appEnvs.JWTAccessExpIn, appEnvs.JWTRefreshExpIn, appEnvs.JWTRefreshRememberMeExpIn, appEnvs.JWTIssuer)
	sessionSrv := sessionService.NewSessionSrv(sessionRepo, authSrv)
	authMethodSrv := authMethodService.NewAuthMethodSrv(userRepo, userFileStg, sessionSrv)
	userSrv := userService.NewUserSrv(userRepo, userFileStg)

	// Handler
	// Note: all subsequent handlers should also be registered using v1
	authMethodHandler.NewAuthMethodHandler(v1, authMethodSrv)
	userHandler.NewUserHandler(v1, userSrv)

	// Initializers
	roleItz := roleInitializer.NewRoleItz(roleRepo)
	if err := roleItz.SeedEssentialRoles(context.Background()); err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := struct {
			Msg string `json:"msg"`
		}{
			Msg: "pong",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	// Templates
	// Redirects requests from "/admin" to the admin home page under API v1
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/v1/admin/home", http.StatusFound)
	})

	adminHandler.NewAdminHandler(v1, appEnvs.ApiBaseUrl)

	return muxWithCors
}

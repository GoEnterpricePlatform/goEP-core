package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	gmailSmtpClient "github.com/GoEnterpricePlatform/goEP-core/internal/gmail-smtp"
	minioClient "github.com/GoEnterpricePlatform/goEP-core/internal/minio"
	mongoClient "github.com/GoEnterpricePlatform/goEP-core/internal/mongo"
	resendClient "github.com/GoEnterpricePlatform/goEP-core/internal/resend"
	adminService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/service"
	authHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/handler"
	authService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/service"
	gmailsmtp "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/adapter/gmail-smtp"
	resendAdapter "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/adapter/resend"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/port"
	mailerService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/service"
	otpCodeRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/repository/mongo"
	otpCodeService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/service"
	permissionInitializer "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/initializer"
	permissionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/repository/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	roleInitializer "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/initializer"
	roleRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/repository/mongo"
	tokenService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/handler"
	postRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/repository/mongo"
	postService "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/service"
	adminHandler "github.com/GoEnterpricePlatform/goEP-core/web/admin/api/handler"
	postHandler "github.com/GoEnterpricePlatform/goEP-core/web/admin/api/posts/handler"
	adminRenderer "github.com/GoEnterpricePlatform/goEP-core/web/admin/renderer"
	publicHandler "github.com/GoEnterpricePlatform/goEP-core/web/public/api/handler"

	cookieService "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/service"

	sessionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/repository/mongo"
	sessionService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/service"
	userFileStorage "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/file-storage/minio"
	userHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/handler"
	userRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/repository/mongo"
	userService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares/logger"

	middlewareServiceTmpl "github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// define app Stack
	appStack := config.NewAppStack(config.DBMongo, config.MailGmail)

	appEnvs := config.Load(appStack)

	zapLogger := logger.NewHttpLogger(appEnvs.AppEnv)

	corsLogger := middlewares.CorsMiddleware(appEnvs.AllowedOrigins)

	// Add global middlewares, the order matters (logger → CORS → router)
	apiHandler := zapLogger.Middleware(corsLogger(mux))

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

	// Resend
	var mailerAdt port.MailerAdt
	switch appStack.Mail {
	case config.MailResend:
		resendCli := resendClient.NewResendClient(appEnvs.ResendApiKey)
		mailerAdt = resendAdapter.NewResendAdt(resendCli, appEnvs.ResendEmailFrom)
	case config.MailGmail:
		gmailSmtp := gmailSmtpClient.NewGmailSmtpClient(appEnvs.GmailUsername, appEnvs.GmailPass, appEnvs.GmailHost)
		mailerAdt = gmailsmtp.NewGmailSmtpAdt(gmailSmtp, appEnvs.GmailAddr, appEnvs.GmailFrom)
	}

	// Collections
	userColl := mongoDB.Collection("users")
	sessionColl := mongoDB.Collection("sessions")
	otpCodeColl := mongoDB.Collection("opt-codes")
	roleColl := mongoDB.Collection("roles")
	permissionColl := mongoDB.Collection("permissions")
	postColl := mongoDB.Collection("posts")

	// Repositories
	userRepo := userRepository.NewUserRepo(mongoConn.DB, userColl)
	sessionRepo := sessionRepository.NewSessionRepo(mongoConn.DB, sessionColl)
	otpCodeRepo := otpCodeRepository.NewOtpCodeRepo(mongoConn.DB, otpCodeColl)
	roleRepo := roleRepository.NewRoleRepo(mongoConn.DB, roleColl)
	permissionRepo := permissionRepository.NewPermissionRepo(mongoConn.DB, permissionColl)
	postRepo := postRepository.NewPostRepo(mongoConn.DB, postColl)

	// Indexes
	err = userRepo.CreateIndexes()
	if err != nil {
		log.Fatal(err)
	}
	err = postRepo.CreateIndexes()
	if err != nil {
		log.Fatal(err)
	}

	// Before indexes, create initializers
	permissionItz := permissionInitializer.NewPermissionItz(permissionRepo)
	permissions, err := permissionItz.SeedEssentialPermissions(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	roleItz := roleInitializer.NewRoleItz(roleRepo)
	if err := roleItz.SeedEssentialRoles(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := roleItz.AddPermissionsToRole(context.Background(), string(domain.RoleSystemAdmin), permissions); err != nil {
		log.Fatal(err)
	}

	// File Storage
	userFileStg := userFileStorage.NewUserFileStg(minioC.Client, appEnvs.MinioBucketName, 0)

	// Services
	tokenSrv := tokenService.NewTokenSrv(appEnvs.JWTAccessSecret, appEnvs.JWTRefreshSecret, appEnvs.JWTAccessExpIn, appEnvs.JWTRefreshExpIn, appEnvs.JWTRefreshRememberMeExpIn, appEnvs.JWTIssuer)
	cookieSrv := cookieService.NewCookieSrv(appEnvs.AppEnv, appEnvs.JwtAccessCookieDur)
	sessionSrv := sessionService.NewSessionSrv(sessionRepo, tokenSrv)
	otpCodeSrv := otpCodeService.NewOtpCodeSrv(otpCodeRepo)
	mailerSrv := mailerService.NewMailerSrv(mailerAdt, appEnvs.AppName)
	authSrv := authService.NewAuthSrv(userRepo, roleRepo, permissionRepo, userFileStg, sessionSrv, otpCodeSrv, mailerSrv)
	userSrv := userService.NewUserSrv(userRepo, userFileStg)
	postSrv := postService.NewPostSrv(postRepo)

	// service - admin
	adminSrv := adminService.NewAdminSrv(userRepo, roleRepo, permissionRepo, sessionSrv)

	// Handler
	// Note: all subsequent handlers should also be registered using v1
	authHandler.NewAuthHandler(v1, authSrv, tokenSrv, appEnvs.AppEnv)
	userHandler.NewUserHandler(v1, userSrv)
	handler.NewPostApiHandler(v1, postSrv)

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
	adminR := adminRenderer.NewAdminRenderer()
	mdwSrvTmpl := middlewareServiceTmpl.NewMdwSrvTmpl(tokenSrv, authSrv, cookieSrv)

	// Templates - admin
	adminH := adminHandler.NewAdminHandler(adminSrv, cookieSrv, appEnvs.ApiBaseUrl, adminR, mdwSrvTmpl)
	adminH.RegisterRoutes(mux, v1)

	// Templates - post
	postH := postHandler.NewPostWebHandler(postSrv, adminH.ApiBaseUrl, adminR)
	postH.RegisterRoutes(v1)

	// Templates - public
	publicHandler.NewPublicHandler(mux, appEnvs.ApiBaseUrl, postSrv)

	return apiHandler
}

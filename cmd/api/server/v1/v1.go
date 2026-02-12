package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/internal/config"
	gmailSmtpClient "github.com/amorindev/go-cms-tmpl/internal/gmail-smtp"
	minioClient "github.com/amorindev/go-cms-tmpl/internal/minio"
	mongoClient "github.com/amorindev/go-cms-tmpl/internal/mongo"
	resendClient "github.com/amorindev/go-cms-tmpl/internal/resend"
	tokenService "github.com/amorindev/go-cms-tmpl/internal/tokens/service"
	adminService "github.com/amorindev/go-cms-tmpl/pkg/identity/admin/service"
	authHandler "github.com/amorindev/go-cms-tmpl/pkg/identity/auth/handler"
	authService "github.com/amorindev/go-cms-tmpl/pkg/identity/auth/service"
	gmailsmtp "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/adapter/gmail-smtp"
	resendAdapter "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/adapter/resend"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/port"
	mailerService "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/service"
	otpCodeRepository "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/repository/mongo"
	otpCodeService "github.com/amorindev/go-cms-tmpl/pkg/identity/opt-codes/service"
	permissionInitializer "github.com/amorindev/go-cms-tmpl/pkg/identity/permissions/initializer"
	permissionRepository "github.com/amorindev/go-cms-tmpl/pkg/identity/permissions/repository/mongo"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/roles/domain"
	roleInitializer "github.com/amorindev/go-cms-tmpl/pkg/identity/roles/initializer"
	roleRepository "github.com/amorindev/go-cms-tmpl/pkg/identity/roles/repository/mongo"
	adminHandler "github.com/amorindev/go-cms-tmpl/web/admin/api/handler"
	adminRenderer "github.com/amorindev/go-cms-tmpl/web/admin/renderer"
	publicHandler "github.com/amorindev/go-cms-tmpl/web/public/api/handler"

	cookieService "github.com/amorindev/go-cms-tmpl/pkg/shared/api/handler/cookie/service"

	sessionRepository "github.com/amorindev/go-cms-tmpl/pkg/identity/session/repository/mongo"
	sessionService "github.com/amorindev/go-cms-tmpl/pkg/identity/session/service"
	userFileStorage "github.com/amorindev/go-cms-tmpl/pkg/identity/users/file-storage/minio"
	userHandler "github.com/amorindev/go-cms-tmpl/pkg/identity/users/handler"
	userRepository "github.com/amorindev/go-cms-tmpl/pkg/identity/users/repository/mongo"
	userService "github.com/amorindev/go-cms-tmpl/pkg/identity/users/service"
	"github.com/amorindev/go-cms-tmpl/pkg/shared/api/middlewares"
	"github.com/amorindev/go-cms-tmpl/pkg/shared/api/middlewares/logger"
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

	// Repositories
	userRepo := userRepository.NewUserRepo(mongoConn.DB, userColl)
	sessionRepo := sessionRepository.NewSessionRepo(mongoConn.DB, sessionColl)
	otpCodeRepo := otpCodeRepository.NewOtpCodeRepo(mongoConn.DB, otpCodeColl)
	roleRepo := roleRepository.NewRoleRepo(mongoConn.DB, roleColl)
	permissionRepo := permissionRepository.NewPermissionRepo(mongoConn.DB, permissionColl)

	// Indexes
	err = userRepo.CreateIndexes()
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
	cookieSrv := cookieService.NewCookieSrv(appEnvs.AppEnv)
	sessionSrv := sessionService.NewSessionSrv(sessionRepo, tokenSrv)
	otpCodeSrv := otpCodeService.NewOtpCodeSrv(otpCodeRepo)
	mailerSrv := mailerService.NewMailerSrv(mailerAdt, appEnvs.AppName)
	authSrv := authService.NewAuthSrv(userRepo, userFileStg, sessionSrv, otpCodeSrv, mailerSrv)
	userSrv := userService.NewUserSrv(userRepo, userFileStg)

	// service - admin
	adminSrv := adminService.NewAdminSrv(userRepo, roleRepo, sessionSrv)

	// Handler
	// Note: all subsequent handlers should also be registered using v1
	authHandler.NewAuthHandler(v1, authSrv, tokenSrv, appEnvs.AppEnv)
	userHandler.NewUserHandler(v1, userSrv)

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
	adminH := adminHandler.NewAdminHandler(adminSrv, cookieSrv, appEnvs.ApiBaseUrl, adminR)
	adminH.RegisterRoutes(mux, v1)

	publicHandler.NewPublicHandler(mux, appEnvs.ApiBaseUrl)

	return apiHandler
}

package config

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	adminService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/service"
	authHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/handler"
	authService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/service"
	gmailsmtp "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/adapter/gmail-smtp"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/adapter/resend"
	mailerService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/service"
	otpCodeRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/repository/mongo"
	otpCodeService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/service"
	permissionInitializer "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/initializer"
	permissionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/repository/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	roleInitializer "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/initializer"
	roleRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/repository/mongo"
	roleService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/service"
	sessionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/repository/mongo"
	sessionService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/service"
	tokenService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/file-storage/disabled"
	userFileStorage "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/file-storage/minio"
	userHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/handler"
	userRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/repository/mongo"
	userService "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/service"
	cookieService "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"

	adminP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/port"
	authP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/port"
	mailerP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/port"
	roleP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/port"
	userFSP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/port"
	userP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/port"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/port"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB *mongo.Database
}

type Module struct {
	AuthSrv    authP.AuthSrv
	UserSrv    userP.UserSrv
	AdminSrv   adminP.AdminSrv
	TokenSrv   tokenP.TokenSrv
	CookieSrv  cookieP.CookieSrv
	AuthApiMdw *middlewares.AuthMiddleware
	RoleSrv    roleP.RoleSrv
}

func NewIdentityModule(cfg ModuleConfig) (*Module, error) {
	// Select mailer provider
	var mailerAdt mailerP.MailerAdt
	switch cfg.AppEnvs.MailerProvider {
	case config.MailResend:
		mailerAdt = resend.NewResendAdt(cfg.AppClients.ResendCli, cfg.AppEnvs.ResendEmailFrom)
	case config.MailGmail:
		mailerAdt = gmailsmtp.NewGmailSmtpAdt(cfg.AppClients.GmailSmtp, cfg.AppEnvs.GmailAddr, cfg.AppEnvs.GmailFrom)
	}

	// File Storage
	var userFileStg userFSP.UserFileStg
	switch cfg.AppEnvs.FileStorageProvider {
	case config.FSMinio:
		userFileStg = userFileStorage.NewUserFileStg(cfg.AppClients.MinioCli.Client, cfg.AppEnvs.MinioBucketName, 0)
	case config.FSOptional:
		userFileStg = disabled.NewDisabledAdapter()
	}

	// module name
	mdlName := "identity"

	// collections
	userCollName := fmt.Sprintf("%s_users", mdlName)
	userColl := cfg.DB.Collection(userCollName)

	sessionCollName := fmt.Sprintf("%s_sessions", mdlName)
	sessionColl := cfg.DB.Collection(sessionCollName)

	otpCodeCollName := fmt.Sprintf("%s_opt-codes", mdlName)
	otpCodeColl := cfg.DB.Collection(otpCodeCollName)

	roleCollName := fmt.Sprintf("%s_roles", mdlName)
	roleColl := cfg.DB.Collection(roleCollName)

	permissionCollName := fmt.Sprintf("%s_permissions", mdlName)
	permissionColl := cfg.DB.Collection(permissionCollName)

	userRepo := userRepository.NewUserRepo(cfg.AppClients.MongoConn.DB, userColl, roleCollName)
	sessionRepo := sessionRepository.NewSessionRepo(cfg.AppClients.MongoConn.DB, sessionColl)
	otpCodeRepo := otpCodeRepository.NewOtpCodeRepo(cfg.AppClients.MongoConn.DB, otpCodeColl)
	roleRepo := roleRepository.NewRoleRepo(cfg.AppClients.MongoConn.DB, roleColl)
	permissionRepo := permissionRepository.NewPermissionRepo(cfg.AppClients.MongoConn.DB, permissionColl)

	// Indexes
	err := userRepo.CreateIndexes()
	if err != nil {
		return nil, err
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

	// Services
	tokenSrv := tokenService.NewTokenSrv(cfg.AppEnvs.JWTAccessSecret, cfg.AppEnvs.JWTRefreshSecret, cfg.AppEnvs.JWTAccessExpIn, cfg.AppEnvs.JWTRefreshExpIn, cfg.AppEnvs.JWTRefreshRememberMeExpIn, cfg.AppEnvs.JWTIssuer)
	authApiMdw := middlewares.NewAuthMdw(tokenSrv)
	cookieSrv := cookieService.NewCookieSrv(cfg.AppEnvs.AppEnv, cfg.AppEnvs.JwtAccessCookieDur)
	sessionSrv := sessionService.NewSessionSrv(sessionRepo, tokenSrv)
	otpCodeSrv := otpCodeService.NewOtpCodeSrv(otpCodeRepo)
	mailerSrv := mailerService.NewMailerSrv(mailerAdt, cfg.AppEnvs.AppName)
	authSrv := authService.NewAuthSrv(userRepo, roleRepo, permissionRepo, userFileStg, sessionSrv, otpCodeSrv, mailerSrv)
	userSrv := userService.NewUserSrv(userRepo, userFileStg, mdlName)
	roleSrv := roleService.NewRoleSrv(roleRepo)

	// service - admin
	adminSrv := adminService.NewAdminSrv(userRepo, roleRepo, permissionRepo, sessionSrv)

	// Register handlers
	authHandler.NewAuthHandler(cfg.APIv1, authSrv, tokenSrv, cfg.AppEnvs.AppEnv, authApiMdw)
	userHandler.NewUserHandler(cfg.APIv1, userSrv, authApiMdw)

	return &Module{
		AuthSrv:    authSrv,
		UserSrv:    userSrv,
		AdminSrv:   adminSrv,
		TokenSrv:   tokenSrv,
		CookieSrv:  cookieSrv,
		AuthApiMdw: authApiMdw,
		RoleSrv:    roleSrv,
	}, nil
}

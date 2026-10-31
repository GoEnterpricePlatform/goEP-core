package config

import (
	"context"
	"fmt"
	"log"
	"net/http"

	adminService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/service"
	authHandler "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/handler"
	authService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/service"
	gmailsmtp "github.com/GoEnterpricePlatform/goEP-core/addons/identity/mailer/adapter/gmail-smtp"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/mailer/adapter/resend"
	mailerService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/mailer/service"
	otpCodeRepository "github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/repository/mongo"
	otpCodeService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/service"
	permissionInitializer "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/initializer"
	permissionRepository "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/repository/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	roleInitializer "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/initializer"
	roleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/repository/mongo"
	roleService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/service"
	sessionRepository "github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/repository/mongo"
	sessionService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/service"
	tokenService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/service"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/file-storage/disabled"
	userFileStorage "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/file-storage/minio"
	userHandler "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/handler"
	userRepository "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/repository/mongo"
	userService "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/service"
	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	cookieService "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/handler/cookie/service"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares"
	"github.com/jackc/pgx/v5"

	adminP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/admin/port"
	authP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/auth/port"
	mailerP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/mailer/port"
	otpCodeP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/port"
	permissionP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/port"
	roleP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	sessionP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"
	userP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/port"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/handler/cookie/port"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB     *mongo.Database
	PGConn *pgx.Conn
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
	var userFileStg userP.UserFileStg
	switch cfg.AppEnvs.FileStorageProvider {
	case config.FSMinio:
		userFileStg = userFileStorage.NewUserFileStg(cfg.AppClients.MinioCli.Client, cfg.AppEnvs.MinioBucketName, 0)
	case config.FSOptional:
		userFileStg = disabled.NewDisabledAdapter()
	}

	// module name
	mdlName := "identity"

	// Database
	var userRepo userP.UserRepo
	var sessionRepo sessionP.SessionRepo
	var otpCodeRepo otpCodeP.OtpCodeRepo
	var roleRepo roleP.RoleRepo
	var permissionRepo permissionP.PermissionRepo

	fmt.Println("---------------------------------")
	switch cfg.AppEnvs.DBProvider {
	case config.DBMongo:
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

		userRepo = userRepository.NewUserRepo(cfg.AppClients.MongoConn.DB, userColl, roleCollName)
		sessionRepo = sessionRepository.NewSessionRepo(cfg.AppClients.MongoConn.DB, sessionColl)
		otpCodeRepo = otpCodeRepository.NewOtpCodeRepo(cfg.AppClients.MongoConn.DB, otpCodeColl)
		roleRepo = roleRepository.NewRoleRepo(cfg.AppClients.MongoConn.DB, roleColl)
		permissionRepo = permissionRepository.NewPermissionRepo(cfg.AppClients.MongoConn.DB, permissionColl)

		// Indexes
		err := userRepo.CreateIndexes()
		if err != nil {
			return nil, err
		}
	case config.DBPostgresql:
	}
	fmt.Println("---------------------------------2")

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
	fmt.Println("---------------------------------3")

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
	fmt.Println("---------------------------------4")

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

package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"

	planPaddleH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/handler"
	planPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	planPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/plans/mongo"
	planPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/plans/postgresql"
	planPaddleService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/service"
	planPaddleTransaction "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/transaction/mongo"
	planPaddlePostgresTransaction "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/transaction/postgresql"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/file-storage/disabled"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/file-storage/minio"
	planH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/handler"
	planP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/port"
	planRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/repository/mongo"
	planPostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/repository/postgresql"
	planService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/service"
	catalogUIHandler "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/ui/handler"
	variationH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/handler"
	varitionP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	varOptionRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/var-option/mongo"
	varOptionPostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/var-option/postgresql"
	variationRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/variation/mongo"
	variationPostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/variation/postgresql"
	variationService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/service"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	frontMiddleware "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux
	TemplateV1 *http.ServeMux
	MdwSrvTmpl *frontMiddleware.MdwSrvTmpl

	DB *mongoDriver.Database

	Deps ModuleDeps
}

// ModuleDeps groups the external dependencies (from other modules/services)
// that this module needs to function.
type ModuleDeps struct {
	AuthApiMdw middlewarePort.Authenticator
}

type Module struct {
	VariationService         varitionP.VariationSrv
	PlanSrv                  planP.PlanSrv
	PaddleCheckoutCatalogSrv planPaddlePort.PaddleCheckoutCatalogSrv
	UI                       *catalogUIHandler.Handler
}

func NewCatalogModule(cfg ModuleConfig) (*Module, error) {
	if cfg.AppEnvs == nil {
		return nil, fmt.Errorf("catalog requires application environment configuration")
	}

	// File Storage
	var planFileStg planP.PlanFileStg
	switch cfg.AppEnvs.FileStorageProvider {
	case config.FSMinio:
		planFileStg = minio.NewPlanFileStg(cfg.AppClients.MinioCli.Client, cfg.AppEnvs.MinioBucketName, 0)
	case config.FSOptional:
		planFileStg = disabled.NewPlanDisabledAdapter()
	}

	// module name
	mdlName := "catalog"

	variationRepo := varitionP.VariationRepo(nil)
	varOptionRepo := varitionP.VarOptionRepo(nil)
	var optionUsage varitionP.OptionUsageChecker
	var planRepo planP.PlanRepo
	var paddleCheckoutCatalogSrv planPaddlePort.PaddleCheckoutCatalogSrv
	var paddlePlanSrv *planPaddleService.Service
	switch cfg.AppEnvs.DBProvider {
	case config.DBMongo:
		if cfg.DB == nil || cfg.AppClients == nil || cfg.AppClients.MongoConn == nil {
			return nil, fmt.Errorf("catalog requires a MongoDB database and client")
		}
		variationsCollName := fmt.Sprintf("%s_variations", mdlName)
		varOptionCollName := fmt.Sprintf("%s_var-options", mdlName)
		variationsColl := cfg.DB.Collection(variationsCollName)
		varOptionColl := cfg.DB.Collection(varOptionCollName)
		plansColl := cfg.DB.Collection(fmt.Sprintf("%s_plans", mdlName))
		variationRepo = variationRepository.NewVariationRepo(cfg.AppClients.MongoConn.DB, variationsColl, varOptionCollName)
		varOptionRepo = varOptionRepository.NewVarOptionRepo(cfg.AppClients.MongoConn.DB, varOptionColl)
		mongoPlanRepo := planRepository.NewPlanRepo(cfg.AppClients.MongoConn.DB, plansColl, varOptionCollName, variationsCollName)
		planRepo = mongoPlanRepo
		optionUsage = mongoPlanRepo
	case config.DBPostgresql:
		if cfg.AppClients == nil || cfg.AppClients.PGConn == nil || cfg.AppClients.PGConn.Pool == nil {
			return nil, fmt.Errorf("catalog requires a PostgreSQL pool")
		}
		pool := cfg.AppClients.PGConn.Pool
		variationRepo = variationPostgresRepository.NewVariationRepo(pool)
		varOptionRepo = varOptionPostgresRepository.NewVarOptionRepo(pool)
		postgresRepo := planPostgresRepository.New(pool)
		planRepo = postgresRepo
		optionUsage = postgresRepo
	default:
		return nil, fmt.Errorf("catalog does not support database provider %q", cfg.AppEnvs.DBProvider)
	}

	// services
	variationSrv := variationService.NewVariationSrv(variationRepo, varOptionRepo, optionUsage)
	planSrv := planService.NewPlanSrv(planRepo, planFileStg, varOptionRepo)

	// register handlers
	variationH.NewVariationHandler(cfg.APIv1, variationSrv, cfg.Deps.AuthApiMdw)
	planH.NewPlanHandler(cfg.APIv1, planSrv, cfg.Deps.AuthApiMdw)
	if cfg.AppEnvs.DBProvider == config.DBMongo {
		plansPaddleColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_plans", mdlName))
		planPaddleRepo := planPaddleRepository.NewPlanPaddleRepo(cfg.AppClients.MongoConn.DB, plansPaddleColl)
		planPaddleTx := planPaddleTransaction.NewPlanPaddleTx(cfg.AppClients.MongoConn.DB, planRepo, planPaddleRepo)
		var paddleSDKClient *paddle.SDK
		if cfg.AppClients.PaddleCli != nil {
			paddleSDKClient = cfg.AppClients.PaddleCli.Client
		}
		paddlePlanSrv = planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, planRepo, planFileStg, paddleSDKClient)
		paddleCheckoutCatalogSrv = paddlePlanSrv
		planPaddleH.NewPlanPaddleHandler(cfg.APIv1, paddlePlanSrv, cfg.Deps.AuthApiMdw, cfg.AppEnvs.PaddleCatalogWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)
	} else if cfg.AppEnvs.DBProvider == config.DBPostgresql {
		pool := cfg.AppClients.PGConn.Pool
		planPaddleRepo := planPaddlePostgresRepository.NewPlanPaddleRepo(pool)
		planPaddleTx := planPaddlePostgresTransaction.NewPlanPaddleTx(pool)
		var paddleSDKClient *paddle.SDK
		if cfg.AppClients.PaddleCli != nil {
			paddleSDKClient = cfg.AppClients.PaddleCli.Client
		}
		paddlePlanSrv = planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, planRepo, planFileStg, paddleSDKClient)
		paddleCheckoutCatalogSrv = paddlePlanSrv
		planPaddleH.NewPlanPaddleHandler(cfg.APIv1, paddlePlanSrv, cfg.Deps.AuthApiMdw, cfg.AppEnvs.PaddleCatalogWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)
	}
	ui := catalogUIHandler.NewHandler(variationSrv, paddlePlanSrv, cfg.MdwSrvTmpl)
	ui.RegisterRoutes(cfg.TemplateV1)

	return &Module{
		VariationService:         variationSrv,
		PlanSrv:                  planSrv,
		PaddleCheckoutCatalogSrv: paddleCheckoutCatalogSrv,
		UI:                       ui,
	}, nil
}

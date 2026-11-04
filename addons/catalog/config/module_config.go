package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"

	planPaddleH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/handler"
	customerPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/customer/mongo"
	customerPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/customer/postgresql"
	planPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/plans/mongo"
	planPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/plans/postgresql"
	subscriptionPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/subscription/mongo"
	subscriptionPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/subscription/postgresql"
	transactionPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/transaction/mongo"
	transactionPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/transaction/postgresql"
	webhookEventPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/webhook-event/mongo"
	webhookEventPaddlePostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/webhook-event/postgresql"
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
	variationH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/handler"
	varitionP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	varOptionRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/var-option/mongo"
	varOptionPostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/var-option/postgresql"
	variationRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/variation/mongo"
	variationPostgresRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/variation/postgresql"
	variationService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/service"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB *mongoDriver.Database

	Deps ModuleDeps
}

// ModuleDeps groups the external dependencies (from other modules/services)
// that this module needs to function.
type ModuleDeps struct {
	AuthApiMdw middlewarePort.Authenticator
}

type Module struct {
	VariationService varitionP.VariationSrv
	PlanSrv          planP.PlanSrv
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
		// Paddle persistence remains Mongo-backed until its PostgreSQL repositories are available.
		plansPaddleColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_plans", mdlName))
		paddleTransactionsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_transactions", mdlName))
		paddleSubscriptionsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_subscriptions", mdlName))
		paddleCustomersColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_customers", mdlName))
		paddleWebhookEventsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_webhook_events", mdlName))
		planPaddleRepo := planPaddleRepository.NewPlanPaddleRepo(cfg.AppClients.MongoConn.DB, plansPaddleColl)
		paddleTransactionRepo := transactionPaddleRepository.NewPaddleTransactionRepo(paddleTransactionsColl)
		paddleSubscriptionRepo := subscriptionPaddleRepository.NewPaddleSubscriptionRepo(paddleSubscriptionsColl)
		paddleCustomerRepo := customerPaddleRepository.NewPaddleCustomerRepo(paddleCustomersColl)
		paddleWebhookEventRepo := webhookEventPaddleRepository.NewPaddleWebhookEventRepo(paddleWebhookEventsColl)
		if err := paddleCustomerRepo.MigrateLegacyIDs(context.Background()); err != nil {
			return nil, fmt.Errorf("migrate Paddle customer IDs: %w", err)
		}
		if err := paddleSubscriptionRepo.MigrateLegacyIDs(context.Background(), paddleCustomerRepo); err != nil {
			return nil, fmt.Errorf("migrate Paddle subscription IDs: %w", err)
		}
		if err := paddleSubscriptionRepo.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := paddleTransactionRepo.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := paddleCustomerRepo.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := paddleWebhookEventRepo.CreateIndexes(); err != nil {
			return nil, err
		}
		planPaddleTx := planPaddleTransaction.NewPlanPaddleTx(cfg.AppClients.MongoConn.DB, planRepo, planPaddleRepo)
		var paddleSDKClient *paddle.SDK
		if cfg.AppClients.PaddleCli != nil {
			paddleSDKClient = cfg.AppClients.PaddleCli.Client
		}
		planPaddleSrv := planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, paddleTransactionRepo, paddleSubscriptionRepo, planRepo, planFileStg, paddleSDKClient, paddleCustomerRepo, paddleWebhookEventRepo)
		planPaddleH.NewPlanPaddleHandler(cfg.APIv1, planPaddleSrv, cfg.Deps.AuthApiMdw, cfg.AppEnvs.PaddleWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)
	} else if cfg.AppEnvs.DBProvider == config.DBPostgresql {
		pool := cfg.AppClients.PGConn.Pool
		planPaddleRepo := planPaddlePostgresRepository.NewPlanPaddleRepo(pool)
		paddleTransactionRepo := transactionPaddlePostgresRepository.NewPaddleTransactionRepo(pool)
		paddleSubscriptionRepo := subscriptionPaddlePostgresRepository.NewPaddleSubscriptionRepo(pool)
		paddleCustomerRepo := customerPaddlePostgresRepository.NewPaddleCustomerRepo(pool)
		paddleWebhookEventRepo := webhookEventPaddlePostgresRepository.NewPaddleWebhookEventRepo(pool)
		if err := paddleWebhookEventRepo.CreateIndexes(); err != nil {
			return nil, err
		}
		planPaddleTx := planPaddlePostgresTransaction.NewPlanPaddleTx(pool)
		var paddleSDKClient *paddle.SDK
		if cfg.AppClients.PaddleCli != nil {
			paddleSDKClient = cfg.AppClients.PaddleCli.Client
		}
		planPaddleSrv := planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, paddleTransactionRepo, paddleSubscriptionRepo, planRepo, planFileStg, paddleSDKClient, paddleCustomerRepo, paddleWebhookEventRepo)
		planPaddleH.NewPlanPaddleHandler(cfg.APIv1, planPaddleSrv, cfg.Deps.AuthApiMdw, cfg.AppEnvs.PaddleWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)
	}

	return &Module{
		VariationService: variationSrv,
		PlanSrv:          planSrv,
	}, nil
}

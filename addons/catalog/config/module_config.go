package config

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"

	planPaddleH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/handler"
	customerPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/customer/mongo"
	planPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/plans/mongo"
	subscriptionPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/subscription/mongo"
	transactionPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/transaction/mongo"
	webhookEventPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/repository/webhook-event/mongo"
	planPaddleService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/service"
	planPaddleTransaction "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/transaction/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/file-storage/disabled"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/file-storage/minio"
	planH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/handler"
	planP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/port"
	planRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/repository/mongo"
	planService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/service"
	variationH "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/handler"
	varitionP "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	varOptionRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/var-option/mongo"
	variationRepository "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/repository/variation/mongo"
	variationService "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/service"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares/port"
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

	// collections
	variationsCollName := fmt.Sprintf("%s_variations", mdlName)
	variationsColl := cfg.DB.Collection(variationsCollName)

	varOptionCollName := fmt.Sprintf("%s_var-options", mdlName)
	varOptionColl := cfg.DB.Collection(varOptionCollName)

	plansCollName := fmt.Sprintf("%s_plans", mdlName)
	plansColl := cfg.DB.Collection(plansCollName)

	plansPaddleCollName := fmt.Sprintf("%s_paddle_plans", mdlName)
	plansPaddleColl := cfg.DB.Collection(plansPaddleCollName)
	paddleTransactionsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_transactions", mdlName))
	paddleSubscriptionsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_subscriptions", mdlName))
	paddleCustomersColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_customers", mdlName))
	paddleWebhookEventsColl := cfg.DB.Collection(fmt.Sprintf("%s_paddle_webhook_events", mdlName))

	// Repositories
	variationRepo := variationRepository.NewVariationRepo(cfg.AppClients.MongoConn.DB, variationsColl, varOptionCollName)
	varOptionRepo := varOptionRepository.NewVarOptionRepo(cfg.AppClients.MongoConn.DB, varOptionColl)
	planRepo := planRepository.NewPlanRepo(cfg.AppClients.MongoConn.DB, plansColl, varOptionCollName, variationsCollName)
	planPaddleRepo := planPaddleRepository.NewPlanPaddleRepo(cfg.AppClients.MongoConn.DB, plansPaddleColl)
	paddleTransactionRepo := transactionPaddleRepository.NewPaddleTransactionRepo(paddleTransactionsColl)
	paddleSubscriptionRepo := subscriptionPaddleRepository.NewPaddleSubscriptionRepo(paddleSubscriptionsColl)
	paddleCustomerRepo := customerPaddleRepository.NewPaddleCustomerRepo(paddleCustomersColl)
	paddleWebhookEventRepo := webhookEventPaddleRepository.NewPaddleWebhookEventRepo(paddleWebhookEventsColl)

	// Migrate billing documents before creating indexes for the new ID fields.
	if err := paddleCustomerRepo.MigrateLegacyIDs(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := paddleSubscriptionRepo.MigrateLegacyIDs(context.Background(), paddleCustomerRepo); err != nil {
		log.Fatal(err)
	}
	// Indexes
	err := paddleSubscriptionRepo.CreateIndexes()
	if err != nil {
		log.Fatal(err)
	}
	if err := paddleTransactionRepo.CreateIndexes(); err != nil {
		log.Fatal(err)
	}
	if err := paddleCustomerRepo.CreateIndexes(); err != nil {
		log.Fatal(err)
	}
	if err := paddleWebhookEventRepo.CreateIndexes(); err != nil {
		log.Fatal(err)
	}

	// transactions
	planPaddleTx := planPaddleTransaction.NewPlanPaddleTx(cfg.AppClients.MongoConn.DB, planRepo, planPaddleRepo)

	// services
	variationSrv := variationService.NewVariationSrv(variationRepo, varOptionRepo, planRepo)
	planSrv := planService.NewPlanSrv(planRepo, planFileStg, varOptionRepo)
	var paddleSDKClient *paddle.SDK
	if cfg.AppClients.PaddleCli != nil {
		paddleSDKClient = cfg.AppClients.PaddleCli.Client
	}
	planPaddleSrv := planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, paddleTransactionRepo, paddleSubscriptionRepo, planRepo, planFileStg, paddleSDKClient, paddleCustomerRepo, paddleWebhookEventRepo)

	// register handlers
	variationH.NewVariationHandler(cfg.APIv1, variationSrv, cfg.Deps.AuthApiMdw)
	planH.NewPlanHandler(cfg.APIv1, planSrv, cfg.Deps.AuthApiMdw)
	planPaddleH.NewPlanPaddleHandler(cfg.APIv1, planPaddleSrv, cfg.Deps.AuthApiMdw, cfg.AppEnvs.PaddleWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)

	return &Module{
		VariationService: variationSrv,
		PlanSrv:          planSrv,
	}, nil
}

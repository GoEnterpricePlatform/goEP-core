package config

import (
	"context"
	"fmt"
	"net/http"

	paddleBillingHandler "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/handler"
	"github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/port"
	customerMongo "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/customer/mongo"
	customerPostgres "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/customer/postgresql"
	subscriptionMongo "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/subscription/mongo"
	subscriptionPostgres "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/subscription/postgresql"
	transactionMongo "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/transaction/mongo"
	transactionPostgres "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/transaction/postgresql"
	webhookMongo "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/webhook-event/mongo"
	webhookPostgres "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/repository/webhook-event/postgresql"
	paddleBillingService "github.com/GoEnterpricePlatform/goEP-core/addons/billing/plans-paddle/service"
	uiHandler "github.com/GoEnterpricePlatform/goEP-core/addons/billing/ui/handler"
	catalogPaddlePort "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	middlewarePort "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares/port"
	frontMiddleware "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs            *appConfig.AppEnvs
	AppClients         *appConfig.AppClients
	APIv1              *http.ServeMux
	AuthApiMdw         middlewarePort.Authenticator
	CatalogCheckoutSrv catalogPaddlePort.PaddleCheckoutCatalogSrv
	TemplateV1         *http.ServeMux
	DB                 *mongo.Database
	PGConn             *pgx.Conn
	MdwSrvTmpl         *frontMiddleware.MdwSrvTmpl
}

type Module struct {
	UI               *uiHandler.Handler
	PaddleBillingSrv port.PaddleBillingSrv
}

func NewBillingModule(cfg ModuleConfig) (*Module, error) {
	if cfg.AppEnvs == nil {
		return nil, fmt.Errorf("billing requires application environment configuration")
	}
	ui := uiHandler.NewHandler(cfg.MdwSrvTmpl)
	ui.RegisterRoutes(cfg.TemplateV1)
	module := &Module{UI: ui}
	switch cfg.AppEnvs.DBProvider {
	case appConfig.DBMongo:
		if cfg.DB == nil {
			return nil, fmt.Errorf("billing requires a MongoDB database")
		}
		transactions := transactionMongo.NewPaddleTransactionRepo(cfg.DB.Collection("billing_paddle_transactions"))
		subscriptions := subscriptionMongo.NewPaddleSubscriptionRepo(cfg.DB.Collection("billing_paddle_subscriptions"))
		customers := customerMongo.NewPaddleCustomerRepo(cfg.DB.Collection("billing_paddle_customers"))
		webhookEvents := webhookMongo.NewPaddleWebhookEventRepo(cfg.DB.Collection("billing_paddle_webhook_events"))
		if err := customers.MigrateLegacyIDs(context.Background()); err != nil {
			return nil, fmt.Errorf("migrate Paddle customer IDs: %w", err)
		}
		if err := subscriptions.MigrateLegacyIDs(context.Background(), customers); err != nil {
			return nil, fmt.Errorf("migrate Paddle subscription IDs: %w", err)
		}
		if err := subscriptions.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := transactions.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := customers.CreateIndexes(); err != nil {
			return nil, err
		}
		if err := webhookEvents.CreateIndexes(); err != nil {
			return nil, err
		}
		var paddleClient *paddleSDK.SDK
		if cfg.AppClients != nil && cfg.AppClients.PaddleCli != nil {
			paddleClient = cfg.AppClients.PaddleCli.Client
		}
		module.PaddleBillingSrv = paddleBillingService.NewPaddleBillingSrv(transactions, subscriptions, customers, webhookEvents, cfg.CatalogCheckoutSrv, paddleClient)
	case appConfig.DBPostgresql:
		if cfg.AppClients == nil || cfg.AppClients.PGConn == nil || cfg.AppClients.PGConn.Pool == nil {
			return nil, fmt.Errorf("billing requires a PostgreSQL pool")
		}
		pool := cfg.AppClients.PGConn.Pool
		transactions := transactionPostgres.NewPaddleTransactionRepo(pool)
		subscriptions := subscriptionPostgres.NewPaddleSubscriptionRepo(pool)
		customers := customerPostgres.NewPaddleCustomerRepo(pool)
		webhookEvents := webhookPostgres.NewPaddleWebhookEventRepo(pool)
		if err := webhookEvents.CreateIndexes(); err != nil {
			return nil, fmt.Errorf("initialize Billing Paddle webhook indexes (apply PostgreSQL migrations with `make migrate-up` first): %w", err)
		}
		var paddleClient *paddleSDK.SDK
		if cfg.AppClients != nil && cfg.AppClients.PaddleCli != nil {
			paddleClient = cfg.AppClients.PaddleCli.Client
		}
		module.PaddleBillingSrv = paddleBillingService.NewPaddleBillingSrv(transactions, subscriptions, customers, webhookEvents, cfg.CatalogCheckoutSrv, paddleClient)
	default:
		return nil, fmt.Errorf("billing does not support database provider %q", cfg.AppEnvs.DBProvider)
	}
	paddleBillingHandler.NewPaddleBillingHandler(cfg.APIv1, module.PaddleBillingSrv, cfg.AuthApiMdw, cfg.AppEnvs.PaddleBillingWebhookSecret, cfg.AppEnvs.IsEnableMorPaddle)
	return module, nil
}

package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"

	planPaddleH "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/handler"
	planPaddleRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/repository/mongo"
	planPaddleService "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/service"
	planPaddleTransaction "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/transaction/mongo"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/file-storage/disabled"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/file-storage/minio"
	planH "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/handler"
	planP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	planRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/repository/mongo"
	planService "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/service"
	variationH "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/handler"
	varitionP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
	varOptionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/repository/var-option/mongo"
	variationRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/repository/variation/mongo"
	variationService "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
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
	AuthApiMdw *middlewares.AuthMiddleware
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

	plansPaddleCollName := fmt.Sprintf("%s_plans_paddle", mdlName)
	plansPaddleColl := cfg.DB.Collection(plansPaddleCollName)

	// Repositories
	variationRepo := variationRepository.NewVariationRepo(cfg.AppClients.MongoConn.DB, variationsColl, varOptionCollName)
	varOptionRepo := varOptionRepository.NewVarOptionRepo(cfg.AppClients.MongoConn.DB, varOptionColl)
	planRepo := planRepository.NewPlanRepo(cfg.AppClients.MongoConn.DB, plansColl, varOptionCollName, variationsCollName)
	planPaddleRepo := planPaddleRepository.NewPlanPaddleRepo(cfg.AppClients.MongoConn.DB, plansPaddleColl)

	// transactions
	planPaddleTx := planPaddleTransaction.NewPlanPaddleTx(cfg.AppClients.MongoConn.DB, planRepo, planPaddleRepo)

	// services
	variationSrv := variationService.NewVariationSrv(variationRepo, varOptionRepo, planRepo)
	planSrv := planService.NewPlanSrv(planRepo, planFileStg, varOptionRepo)
	planPaddleSrv := planPaddleService.NewPlanPaddleSrv(planPaddleTx, planPaddleRepo, planRepo, planFileStg)

	// register handlers
	variationH.NewVariationHandler(cfg.APIv1, variationSrv, cfg.Deps.AuthApiMdw)
	planH.NewPlanHandler(cfg.APIv1, planSrv, cfg.Deps.AuthApiMdw)
	planPaddleH.NewPlanPaddleHandler(cfg.APIv1, planPaddleSrv, cfg.Deps.AuthApiMdw)

	return &Module{
		VariationService: variationSrv,
		PlanSrv:          planSrv,
	}, nil
}

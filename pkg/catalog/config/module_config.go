package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"

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

	// Repositories
	variationRepo := variationRepository.NewVariationRepo(cfg.AppClients.MongoConn.DB, variationsColl)
	varOptionRepo := varOptionRepository.NewVarOptionRepo(cfg.AppClients.MongoConn.DB, varOptionColl)
	planRepo := planRepository.NewPlanRepo(cfg.AppClients.MongoConn.DB, plansColl)

	// services
	variationSrv := variationService.NewVariationSrv(variationRepo, varOptionRepo)
	planSrv := planService.NewPlanSrv(planRepo, planFileStg, varOptionRepo)

	// register handlers
	variationH.NewVariationHandler(cfg.APIv1, variationSrv, cfg.Deps.AuthApiMdw)
	planH.NewPlanHandler(cfg.APIv1, planSrv, cfg.Deps.AuthApiMdw)

	return &Module{
		VariationService: variationSrv,
		PlanSrv:          planSrv,
	}, nil
}

package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/handler"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"
	varOptionRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/repository/var-option/mongo"
	variationRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/repository/variation/mongo"
	variationService "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppStack   *config.AppStack
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB *mongo.Database
}

type Module struct {
	VariationService port.VariationSrv
}

func NewCatalogModule(cfg ModuleConfig) (*Module, error) {

	// module name
	mdlName := "catalog"

	// collections
	variationsCollName := fmt.Sprintf("%s_variations", mdlName)
	variationsColl := cfg.DB.Collection(variationsCollName)

	varOptionCollName := fmt.Sprintf("%s_var-options", mdlName)
	varOptionColl := cfg.DB.Collection(varOptionCollName)

	// Repositories
	variationRepo := variationRepository.NewVariationRepo(cfg.AppClients.MongoConn.DB, variationsColl)
	varOptionRepo := varOptionRepository.NewVarOptionRepo(cfg.AppClients.MongoConn.DB, varOptionColl)

	// services
	variationSrv := variationService.NewVariationSrv(variationRepo, varOptionRepo)

	// register handlers
	handler.NewVariationHandler(cfg.APIv1, variationSrv)

	return &Module{
		VariationService: variationSrv,
	}, nil
}

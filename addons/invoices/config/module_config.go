package config

import (
	"net/http"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	DB         *mongo.Database
}

type Module struct{}

func NewInvoicesModule(_ ModuleConfig) (*Module, error) {
	return &Module{}, nil
}

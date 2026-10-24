package modules

import (
	"fmt"
	"net/http"
	"reflect"

	appConfig "github.com/GoEnterpricePlatform/goEP-core/internal/config"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BackendRuntime contains shared infrastructure passed to a module-local
// initializer. Module-specific services stay behind typed contracts.

// mas adelante: en ves de pasar *mongo.Database y *pgx.Conn seria pasarle  Data o agrupar
// ambos DBData y poner ahi ambos clientes 
type BackendRuntime struct {
	AppEnvs    *appConfig.AppEnvs
	AppClients *appConfig.AppClients
	APIv1      *http.ServeMux
	DB         *mongo.Database
	PGConn 		*pgx.Conn
	Services   *Services
}

// Services stores module contracts by their Go interface or concrete type.
type Services struct {
	values map[reflect.Type]any
}

// NewServices creates an empty service registry for one backend startup.
func NewServices() *Services {
	return &Services{values: make(map[reflect.Type]any)}
}

// ProvideService publishes a module service using its declared contract type.
func ProvideService[T any](services *Services, value T) error {
	if services == nil {
		return fmt.Errorf("module services are nil")
	}
	typeOf := reflect.TypeOf((*T)(nil)).Elem()
	if _, exists := services.values[typeOf]; exists {
		return fmt.Errorf("module service %q is already provided", typeOf)
	}
	services.values[typeOf] = value
	return nil
}

// RequireService retrieves a previously provided service by its typed contract.
func RequireService[T any](services *Services) (T, error) {
	var zero T
	if services == nil {
		return zero, fmt.Errorf("module services are nil")
	}
	typeOf := reflect.TypeOf((*T)(nil)).Elem()
	value, exists := services.values[typeOf]
	if !exists {
		return zero, fmt.Errorf("required module service %q is not available", typeOf)
	}
	service, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("provided module service %q has an incompatible type", typeOf)
	}
	return service, nil
}

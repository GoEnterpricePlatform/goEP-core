package config

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	pProviderRepo "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/repository/mongo"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/keystore"

	paymentProviderHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/handler"
	paymentProviderService "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/service"

	pProviderP "github.com/GoEnterpricePlatform/goEP-core/pkg/billing/payment-providers/port"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/encryptor"

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
	PProviderService pProviderP.PaymentProviderSrv
	EncryptorService encryptor.EncryptorSrv
}

func NewBillingModule(cfg ModuleConfig) (*Module, error) {

	// module name
	mdlName := "billing"

	// collections
	collenctionName := fmt.Sprintf("%s_payment-providers", mdlName)
	pProviderColl := cfg.DB.Collection(collenctionName)

	// Repositories
	pProviderRepo := pProviderRepo.NewPaymentProviderRepo(cfg.AppClients.MongoConn.DB, pProviderColl)

	// Indexes
	err := pProviderRepo.CreateIndexes()
	if err != nil {
		log.Fatal(err)
	}

	// Services - billing
	// For now, FileKeyStore depends on keyFile = "payment-encryption.key", the FileKeyStore
	// should be dynamic, for example pass it the name of the key my-key and inside it chains
	// it with .key also to save it inside goep_core_data/ it should be the name of the module
	//  followed by the folder keys/ followed by the key, for example,
	// goep_core_data/mi-name-module/keys/my-key.key
	keyStore := keystore.NewFileKeyStore(cfg.AppEnvs.AppDataPath)

	key, err := keyStore.LoadOrCreateKey()
	if err != nil {
		log.Fatal(err)
	}

	encryptorSrv := encryptor.NewAESGCM(key)
	pProviderSrv := paymentProviderService.NewPaymentProviderSrv(pProviderRepo, encryptorSrv)

	// register handlers
	paymentProviderHandler.NewPaymentProviderApiHandler(cfg.APIv1, pProviderSrv)

	return &Module{
		PProviderService: pProviderSrv,
		EncryptorService: encryptorSrv,
	}, nil
}

package mongo

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ port.PaddleTransactionRepo = &TransactionRepository{}

type TransactionRepository struct {
	Collection *mongo.Collection
}

func NewPaddleTransactionRepo(collection *mongo.Collection) *TransactionRepository {
	return &TransactionRepository{Collection: collection}
}

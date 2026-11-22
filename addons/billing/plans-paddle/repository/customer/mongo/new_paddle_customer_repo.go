package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewPaddleCustomerRepo(collection *mongo.Collection) *CustomerRepository {
	return &CustomerRepository{Collection: collection}
}

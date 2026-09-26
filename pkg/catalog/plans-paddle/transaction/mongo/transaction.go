package mongo

import (
	paddlePlanP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ paddlePlanP.PaddlePlanTx = &Transaction{}

type Transaction struct {
	Client         *mongo.Client
	PlanRepo       port.PlanRepo
	PaddlePlanRepo paddlePlanP.PaddlePlanRepo
}

func NewPlanPaddleTx(client *mongo.Client, planRepo port.PlanRepo, paddlePlanRepo paddlePlanP.PaddlePlanRepo) *Transaction {
	return &Transaction{
		Client:         client,
		PlanRepo:       planRepo,
		PaddlePlanRepo: paddlePlanRepo,
	}
}

package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/transaction"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *TransactionRepository) Insert(ctx context.Context, transaction *domain.PaddleTransaction) error {
	if transaction == nil || transaction.PaddleTransactionID == "" {
		return fmt.Errorf("Paddle transaction is missing Paddle transaction ID")
	}
	if transaction.ID == "" {
		transaction.ID = bson.NewObjectID().Hex()
	}
	transactionModel := model.FromDomainPaddleTransaction(transaction)
	if _, err := r.Collection.InsertOne(ctx, transactionModel); err != nil {
		return fmt.Errorf("error saving Paddle transaction: %w", err)
	}
	return nil
}

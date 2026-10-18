package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/transaction"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *TransactionRepository) FindAllLinkedToSubscriptions(ctx context.Context) ([]*domain.PaddleTransaction, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"subscription_id": bson.M{"$type": "string", "$ne": ""}})
	if err != nil {
		return nil, fmt.Errorf("error finding transactions linked to subscriptions: %w", err)
	}
	defer cursor.Close(ctx)
	var docs []model.PaddleTransactionModel
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("error decoding transactions linked to subscriptions: %w", err)
	}
	transactions := make([]*domain.PaddleTransaction, 0, len(docs))
	for i := range docs {
		transactions = append(transactions, docs[i].ToDomain())
	}
	return transactions, nil
}

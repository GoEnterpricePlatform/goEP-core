package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/transaction"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *TransactionRepository) FindAllByUserID(ctx context.Context, userID string) ([]*domain.PaddleTransaction, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("error finding Paddle transactions for user: %w", err)
	}
	defer cursor.Close(ctx)
	var docs []model.PaddleTransactionModel
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("error decoding Paddle transactions for user: %w", err)
	}
	transactions := make([]*domain.PaddleTransaction, 0, len(docs))
	for i := range docs {
		transactions = append(transactions, docs[i].ToDomain())
	}
	return transactions, nil
}

package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *TransactionRepository) FindByPaddleTransactionID(ctx context.Context, paddleTransactionID string) (*domain.PaddleTransaction, error) {
	var doc model.PaddleTransactionModel
	if err := r.Collection.FindOne(ctx, bson.M{"paddle_transaction_id": paddleTransactionID}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle transaction: %w", err)
	}
	return doc.ToDomain(), nil
}

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

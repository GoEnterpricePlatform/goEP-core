package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/transaction"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

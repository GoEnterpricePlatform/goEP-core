package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type legacyPaddleCheckoutModel struct {
	PaddleTransactionID string        `bson:"_id"`
	UserID              string        `bson:"user_id,omitempty"`
	SubscriptionID      *string       `bson:"subscription_id,omitempty"`
	PlanID              string        `bson:"plan_id"`
	PlanItemID          string        `bson:"plan_item_id"`
	Status              string        `bson:"status"`
	UpdatedAt           string        `bson:"updated_at"`
	LastEventAt         bson.DateTime `bson:"last_event_at,omitempty"`
}

// MigrateLegacyCheckouts moves checkout documents into the transaction
// collection and replaces external IDs in references with local document IDs.
func (r *TransactionRepository) MigrateLegacyCheckouts(ctx context.Context, legacyCollection *mongo.Collection, subscriptionRepo port.PaddleSubscriptionRepo, customerRepo port.PaddleCustomerRepo) error {
	cursor, err := legacyCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error listing legacy Paddle checkouts: %w", err)
	}
	defer cursor.Close(ctx)
	var legacyTransactions []legacyPaddleCheckoutModel
	if err := cursor.All(ctx, &legacyTransactions); err != nil {
		return fmt.Errorf("error decoding legacy Paddle checkouts: %w", err)
	}
	for _, legacy := range legacyTransactions {
		var existing model.PaddleTransactionModel
		findErr := r.Collection.FindOne(ctx, bson.M{"paddle_transaction_id": legacy.PaddleTransactionID}).Decode(&existing)
		if findErr == nil {
			if _, err := legacyCollection.DeleteOne(ctx, bson.M{"_id": legacy.PaddleTransactionID}); err != nil {
				return fmt.Errorf("error deleting migrated legacy Paddle checkout %s: %w", legacy.PaddleTransactionID, err)
			}
			continue
		}
		if findErr != mongo.ErrNoDocuments {
			return fmt.Errorf("error checking migrated Paddle transaction %s: %w", legacy.PaddleTransactionID, findErr)
		}

		transaction := &model.PaddleTransactionModel{
			ID:                  bson.NewObjectID().Hex(),
			PaddleTransactionID: legacy.PaddleTransactionID,
			UserID:              legacy.UserID,
			PlanID:              legacy.PlanID,
			PlanItemID:          legacy.PlanItemID,
			Status:              legacy.Status,
			UpdatedAt:           legacy.UpdatedAt,
			LastEventAt:         legacy.LastEventAt.Time(),
		}
		if legacy.UserID != "" {
			if customer, err := customerRepo.FindByUserID(ctx, legacy.UserID); err == nil {
				transaction.CustomerID = customer.ID
			} else if err != sharedD.ErrNotFound {
				return fmt.Errorf("error finding customer for migrated Paddle transaction %s: %w", legacy.PaddleTransactionID, err)
			}
		}
		if legacy.SubscriptionID != nil && *legacy.SubscriptionID != "" {
			if subscription, err := subscriptionRepo.FindByPaddleSubscriptionID(ctx, *legacy.SubscriptionID); err == nil {
				transaction.SubscriptionID = &subscription.ID
			} else if err != sharedD.ErrNotFound {
				return fmt.Errorf("error finding subscription for migrated Paddle transaction %s: %w", legacy.PaddleTransactionID, err)
			}
		}
		if _, err := r.Collection.InsertOne(ctx, transaction); err != nil {
			return fmt.Errorf("error inserting migrated Paddle transaction %s: %w", legacy.PaddleTransactionID, err)
		}
		if _, err := legacyCollection.DeleteOne(ctx, bson.M{"_id": legacy.PaddleTransactionID}); err != nil {
			return fmt.Errorf("error deleting migrated legacy Paddle checkout %s: %w", legacy.PaddleTransactionID, err)
		}
	}
	return nil
}

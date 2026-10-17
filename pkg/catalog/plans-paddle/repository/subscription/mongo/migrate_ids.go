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

// MigrateLegacyIDs gives legacy subscriptions independent local IDs and rewrites
// their customer references to the new customer document IDs.
func (r *SubscriptionRepository) MigrateLegacyIDs(ctx context.Context, customerRepo port.PaddleCustomerRepo) error {
	cursor, err := r.Collection.Find(ctx, bson.M{"paddle_subscription_id": bson.M{"$exists": false}})
	if err != nil {
		return fmt.Errorf("error listing legacy Paddle subscriptions: %w", err)
	}
	defer cursor.Close(ctx)
	var legacySubscriptions []model.PaddleSubscriptionModel
	if err := cursor.All(ctx, &legacySubscriptions); err != nil {
		return fmt.Errorf("error decoding legacy Paddle subscriptions: %w", err)
	}
	for _, legacy := range legacySubscriptions {
		if legacy.ID == "" {
			continue
		}
		paddleSubscriptionID := legacy.ID
		var existing model.PaddleSubscriptionModel
		findErr := r.Collection.FindOne(ctx, bson.M{"paddle_subscription_id": paddleSubscriptionID}).Decode(&existing)
		if findErr == nil {
			if _, err := r.Collection.DeleteOne(ctx, bson.M{"_id": legacy.ID}); err != nil {
				return fmt.Errorf("error deleting migrated legacy Paddle subscription %s: %w", paddleSubscriptionID, err)
			}
			continue
		}
		if findErr != mongo.ErrNoDocuments {
			return fmt.Errorf("error checking migrated Paddle subscription %s: %w", paddleSubscriptionID, findErr)
		}
		legacy.ID = bson.NewObjectID().Hex()
		legacy.PaddleSubscriptionID = paddleSubscriptionID
		if customer, err := customerRepo.FindByUserID(ctx, legacy.UserID); err == nil {
			legacy.CustomerID = customer.ID
		} else if err != nil && err != sharedD.ErrNotFound {
			return fmt.Errorf("error finding customer for migrated Paddle subscription %s: %w", paddleSubscriptionID, err)
		} else {
			legacy.CustomerID = ""
		}
		if _, err := r.Collection.InsertOne(ctx, legacy); err != nil {
			return fmt.Errorf("error inserting migrated Paddle subscription %s: %w", paddleSubscriptionID, err)
		}
		if _, err := r.Collection.DeleteOne(ctx, bson.M{"_id": paddleSubscriptionID}); err != nil {
			return fmt.Errorf("error deleting migrated legacy Paddle subscription %s: %w", paddleSubscriptionID, err)
		}
	}
	return nil
}

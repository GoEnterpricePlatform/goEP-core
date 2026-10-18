package mongo

import (
	"context"
	"fmt"

	model "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/customer"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MigrateLegacyIDs assigns independent document IDs to customer documents that
// were previously keyed directly by GoEP user ID.
func (r *CustomerRepository) MigrateLegacyIDs(ctx context.Context) error {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": bson.M{"$exists": false}})
	if err != nil {
		return fmt.Errorf("error listing legacy Paddle customers: %w", err)
	}
	defer cursor.Close(ctx)
	var legacyCustomers []model.PaddleCustomerModel
	if err := cursor.All(ctx, &legacyCustomers); err != nil {
		return fmt.Errorf("error decoding legacy Paddle customers: %w", err)
	}
	for _, legacy := range legacyCustomers {
		if legacy.ID == "" || legacy.PaddleCustomerID == "" {
			continue
		}
		var existing model.PaddleCustomerModel
		findErr := r.Collection.FindOne(ctx, bson.M{"user_id": legacy.ID}).Decode(&existing)
		if findErr == nil {
			if _, err := r.Collection.DeleteOne(ctx, bson.M{"_id": legacy.ID}); err != nil {
				return fmt.Errorf("error deleting migrated legacy Paddle customer %s: %w", legacy.ID, err)
			}
			continue
		}
		if findErr != mongo.ErrNoDocuments {
			return fmt.Errorf("error checking migrated Paddle customer %s: %w", legacy.ID, findErr)
		}
		migrated := &model.PaddleCustomerModel{
			ID:               bson.NewObjectID().Hex(),
			UserID:           legacy.ID,
			PaddleCustomerID: legacy.PaddleCustomerID,
		}
		if _, err := r.Collection.InsertOne(ctx, migrated); err != nil {
			return fmt.Errorf("error inserting migrated Paddle customer %s: %w", legacy.ID, err)
		}
		if _, err := r.Collection.DeleteOne(ctx, bson.M{"_id": legacy.ID}); err != nil {
			return fmt.Errorf("error deleting migrated legacy Paddle customer %s: %w", legacy.ID, err)
		}
	}
	return nil
}

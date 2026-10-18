package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	model "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model/customer"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *CustomerRepository) Upsert(ctx context.Context, customer *domain.PaddleCustomer) error {
	if customer == nil || customer.UserID == "" || customer.PaddleCustomerID == "" {
		return fmt.Errorf("Paddle customer requires user ID and Paddle customer ID")
	}
	doc := model.FromDomainPaddleCustomer(customer)
	var existing model.PaddleCustomerModel
	findErr := r.Collection.FindOne(ctx, bson.M{"user_id": doc.UserID}).Decode(&existing)
	if findErr == nil {
		doc.ID = existing.ID
	} else if findErr != mongo.ErrNoDocuments {
		return fmt.Errorf("error loading Paddle customer before upsert: %w", findErr)
	}
	if doc.ID == "" {
		doc.ID = bson.NewObjectID().Hex()
	}
	_, err := r.Collection.UpdateOne(ctx,
		bson.M{"user_id": doc.UserID},
		bson.M{"$set": bson.M{"paddle_customer_id": doc.PaddleCustomerID}, "$setOnInsert": bson.M{"_id": doc.ID, "user_id": doc.UserID}},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("error upserting Paddle customer: %w", err)
	}
	// Another webhook may have inserted this user between the initial lookup
	// and UpdateOne. In that case $setOnInsert does not use doc.ID, so read the
	// persisted document to return its actual internal ID to related records.
	if err := r.Collection.FindOne(ctx, bson.M{"user_id": doc.UserID}).Decode(&existing); err != nil {
		return fmt.Errorf("error loading Paddle customer after upsert: %w", err)
	}
	customer.ID = existing.ID
	return nil
}

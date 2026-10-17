package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ port.PaddleCustomerRepo = &CustomerRepository{}

type CustomerRepository struct {
	Collection *mongo.Collection
}

func NewPaddleCustomerRepo(collection *mongo.Collection) *CustomerRepository {
	return &CustomerRepository{Collection: collection}
}

func (r *CustomerRepository) Upsert(ctx context.Context, customer *domain.PaddleCustomer) error {
	if customer == nil || customer.UserID == "" || customer.PaddleCustomerID == "" {
		return fmt.Errorf("Paddle customer requires user ID and Paddle customer ID")
	}
	doc := model.FromDomainPaddleCustomer(customer)
	fmt.Println(doc.ID)
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
	// doc.ID
	return nil
}

func (r *CustomerRepository) FindByUserID(ctx context.Context, userID string) (*domain.PaddleCustomer, error) {
	var doc model.PaddleCustomerModel
	if err := r.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle customer by user ID: %w", err)
	}
	return doc.ToDomain(), nil
}

func (r *CustomerRepository) FindByPaddleCustomerID(ctx context.Context, paddleCustomerID string) (*domain.PaddleCustomer, error) {
	var doc model.PaddleCustomerModel
	if err := r.Collection.FindOne(ctx, bson.M{"paddle_customer_id": paddleCustomerID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error finding Paddle customer by Paddle customer ID: %w", err)
	}
	return doc.ToDomain(), nil
}

func (r *CustomerRepository) CreateIndexes() error {
	_, err := r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetName("paddle_customer_user_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle customer user index: %w", err)
	}
	_, err = r.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "paddle_customer_id", Value: 1}},
		Options: options.Index().SetName("paddle_customer_external_id").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating Paddle customer external ID index: %w", err)
	}
	return nil
}

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

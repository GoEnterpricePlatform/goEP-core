package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) Update(ctx context.Context, id string, role *domain.Role) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return sharedD.ErrIncorrectID
	}

	update := bson.M{
		"$set": bson.M{
			"name": role.Name,
		},
	}

	filter := bson.M{"_id": oID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sharedD.ErrNotFound
		}
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

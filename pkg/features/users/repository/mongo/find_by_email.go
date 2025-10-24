package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/features/users/domain"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	filter := bson.D{{Key: "email", Value: email}}

	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: user with email %s not found: %w", dShared.ErrNotFound, email, err)
		}
		return nil, fmt.Errorf("user with email %s not found: %w", email, err)
	}

	oID, ok := user.ID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("ID is not a bson.ObjectID")
	}
	user.ID = oID.Hex()

	return &user, nil
}

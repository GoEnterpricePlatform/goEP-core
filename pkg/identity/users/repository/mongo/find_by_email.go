package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/model"
	dShared "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userModel model.UserNoSqlModel

	filter := bson.D{{Key: "email", Value: email}}

	err := r.Collection.FindOne(ctx, filter).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: user with email %s not found: %w", dShared.ErrNotFound, email, err)
		}
		return nil, fmt.Errorf("user with email %s not found: %w", email, err)
	}

	var user domain.User
	userModel.ToDomain(&user)

	return &user, nil
}


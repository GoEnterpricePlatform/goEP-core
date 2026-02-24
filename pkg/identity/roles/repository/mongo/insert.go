package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/domain"
	sharedDomain "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/roles/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, role *domain.Role) error {
	id := bson.NewObjectID()

	roleModel, err := model.FromDomain(role, id)
	if err != nil {
		return fmt.Errorf("error mapping role to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, roleModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting role: %w", sharedDomain.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting role: %w", err)
	}

	roleModel.ToDomain(role)

	return nil
}

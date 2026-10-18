package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Find(ctx context.Context, id string) (*domain.Role, error) {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	var roleModel model.RoleNoSqlModel
	err = r.Collection.FindOne(ctx, bson.M{"_id": oID}).Decode(&roleModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("error getting role: %w", err)
	}

	var role domain.Role
	roleModel.ToDomain(&role)

	return &role, nil
}

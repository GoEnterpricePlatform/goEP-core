package mongo

import (
	"context"
	"fmt"

	"github.com/amorindev/go-cms-tmpl/pkg/features/permissions/domain"
	"github.com/amorindev/go-cms-tmpl/pkg/features/permissions/model"
	sharedDomain "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, permission *domain.Permission) error {
	id := bson.NewObjectID()

	permissionModel := model.FromDomain(permission, id)

	_, err := r.Collection.InsertOne(ctx, permissionModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting permission: %w", sharedDomain.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting permission: %w", err)
	}
	
	permissionModel.ToDomain(permission)

	return nil
}

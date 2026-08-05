package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/billing/variations/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) Insert(ctx context.Context, varOption *domain.VarOption) error {
	id := bson.NewObjectID()
	
	varOptionModel, err := model.FromDomainOption(varOption, id)
	if err != nil {
		return fmt.Errorf("error mapping varOption to mongo model: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, varOptionModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: error inserting varOption: %s", sharedD.ErrDuplicateKey, err.Error())
		}
		return fmt.Errorf("error inserting varOption: %s", err.Error())
	}

	varOptionModel.ToDomain(varOption)

	return nil
}
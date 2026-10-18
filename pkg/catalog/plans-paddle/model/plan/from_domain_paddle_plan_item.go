package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func FromDomainPaddlePlanItem(o *domain.PaddlePlanItem, id bson.ObjectID) (*PaddlePlanItemModel, error) {
	if o == nil {
		return nil, nil
	}

	planItemID, err := bson.ObjectIDFromHex(o.PlanItemID)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	return &PaddlePlanItemModel{
		ID:            id,
		PlanItemID:    planItemID,
		PaddlePriceID: o.PaddlePriceID,
		Status:        o.Status,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}, nil
}

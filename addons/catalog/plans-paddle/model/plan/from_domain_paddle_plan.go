package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func FromDomainPaddlePlan(d *domain.PaddlePlan, id bson.ObjectID) (*PaddlePlanNoSqlModel, error) {
	if d == nil {
		return nil, nil
	}

	planID, err := bson.ObjectIDFromHex(d.PlanID)
	if err != nil {
		return nil, sharedD.ErrIncorrectID
	}

	items := make([]*PaddlePlanItemModel, 0, len(d.Items))

	for _, item := range d.Items {
		if item == nil {
			continue
		}

		itemID := bson.NewObjectID()
		if item.ID != "" {
			itemID, err = bson.ObjectIDFromHex(item.ID)
			if err != nil {
				return nil, sharedD.ErrIncorrectID
			}
		}

		itemModel, err := FromDomainPaddlePlanItem(item, itemID)
		if err != nil {
			return nil, err
		}

		items = append(items, itemModel)
	}

	return &PaddlePlanNoSqlModel{
		ID:              id,
		PlanID:          planID,
		PaddleProductID: d.PaddleProductID,
		Order:           d.Order,
		Items:           items,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}, nil
}

package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaddlePlanNoSqlModel struct {
	ID              bson.ObjectID          `bson:"_id"`
	PlanID          bson.ObjectID          `bson:"plan_id"`
	PaddleProductID string                 `bson:"paddle_product_id"`
	Items           []*PaddlePlanItemModel `bson:"items"`
	CreatedAt       *time.Time             `bson:"created_at"`
	UpdatedAt       *time.Time             `bson:"updated_at"`
}

type PaddlePlanItemModel struct {
	ID            bson.ObjectID `bson:"_id"`
	PlanItemID    bson.ObjectID `bson:"plan_item_id"`
	PaddlePriceID string        `bson:"paddle_price_id"`
	CreatedAt     *time.Time    `bson:"created_at"`
	UpdatedAt     *time.Time    `bson:"updated_at"`
}

func (m *PaddlePlanItemModel) ToDomain(o *domain.PaddlePlanItem) {
	if m == nil || o == nil {
		return
	}

	o.ID = m.ID.Hex()
	o.PlanItemID = m.PlanItemID.Hex()
	o.PaddlePriceID = m.PaddlePriceID
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

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
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}, nil
}

func (m *PaddlePlanNoSqlModel) ToDomain(o *domain.PaddlePlan) {
	if m == nil || o == nil {
		return
	}

	planID := m.PlanID.Hex()

	items := make([]*domain.PaddlePlanItem, 0, len(m.Items))

	for _, itemModel := range m.Items {
		if itemModel == nil {
			continue
		}

		item := &domain.PaddlePlanItem{}
		itemModel.ToDomain(item)

		items = append(items, item)
	}

	o.ID = m.ID.Hex()
	o.PlanID = planID
	o.PaddleProductID = m.PaddleProductID
	o.Items = items
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

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
		Items:           items,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}, nil
}

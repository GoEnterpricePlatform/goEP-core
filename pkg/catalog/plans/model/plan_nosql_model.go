package model

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PlanNoSqlModel struct {
	ID          bson.ObjectID         `bson:"_id"`
	Name        string                `bson:"name"`
	Description *string               `bson:"description"`
	ImgPath     *string               `bson:"img_path"`
	Items       []*PlanItemNoSqlModel `bson:"items,omitempty"`
	CreatedAt   *time.Time            `bson:"created_at"`
	UpdatedAt   *time.Time            `bson:"updated_at"`
}

func (m *PlanNoSqlModel) ToDomain(v *domain.Plan) {
	if m == nil {
		return
	}

	planItems := make([]*domain.PlanItem, 0, len(m.Items))

	for _, item := range m.Items {
		planItem := &domain.PlanItem{}
		item.ToDomain(planItem)

		planItems = append(planItems, planItem)
	}

	v.ID = m.ID.Hex()
	v.Name = m.Name
	v.Description = m.Description
	v.ImgPath = m.ImgPath
	v.Items = planItems
	v.CreatedAt = m.CreatedAt
	v.UpdatedAt = m.UpdatedAt
}

func FromDomainPlan(d *domain.Plan, id bson.ObjectID) (*PlanNoSqlModel, error) {
	if d == nil {
		return nil, nil
	}

	planItems := make([]*PlanItemNoSqlModel, 0, len(d.Items))

	for _, item := range d.Items {
		if item == nil {
			continue
		}

		itemID := bson.NewObjectID()

		model, err := FromDomainPlanItem(item, itemID)
		if err != nil {
			return nil, err
		}

		planItems = append(planItems, model)
	}

	return &PlanNoSqlModel{
		ID:          id,
		Name:        d.Name,
		Description: d.Description,
		ImgPath:     d.ImgPath,
		Items:       planItems,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

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
	o.Order = m.Order
	o.Items = items
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

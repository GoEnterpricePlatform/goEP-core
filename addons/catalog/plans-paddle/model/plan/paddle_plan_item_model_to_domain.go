package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func (m *PaddlePlanItemModel) ToDomain(o *domain.PaddlePlanItem) {
	if m == nil || o == nil {
		return
	}

	o.ID = m.ID.Hex()
	o.PlanItemID = m.PlanItemID.Hex()
	o.PaddlePriceID = m.PaddlePriceID
	o.Status = m.Status
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

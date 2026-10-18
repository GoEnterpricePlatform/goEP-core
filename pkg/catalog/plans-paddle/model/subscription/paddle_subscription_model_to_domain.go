package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

func (m *PaddleSubscriptionModel) ToDomain() *domain.PaddleSubscription {
	if m == nil {
		return nil
	}
	return &domain.PaddleSubscription{
		ID:                    m.ID,
		PaddleSubscriptionID:  m.PaddleSubscriptionID,
		UserID:                m.UserID,
		CustomerID:            m.CustomerID,
		PlanID:                m.PlanID,
		PlanItemID:            m.PlanItemID,
		Status:                m.Status,
		CurrentPeriodStartsAt: m.CurrentPeriodStartsAt,
		CurrentPeriodEndsAt:   m.CurrentPeriodEndsAt,
		NextBilledAt:          m.NextBilledAt,
		CanceledAt:            m.CanceledAt,
		ScheduledChange:       toDomainSubscriptionChange(m.ScheduledChange),
		UpdatedAt:             m.UpdatedAt,
	}
}

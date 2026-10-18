package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func toDomainSubscriptionChange(change *PaddleSubscriptionChangeModel) *domain.PaddleSubscriptionChange {
	if change == nil {
		return nil
	}
	return &domain.PaddleSubscriptionChange{Action: change.Action, EffectiveAt: change.EffectiveAt}
}

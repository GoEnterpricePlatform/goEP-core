package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
)

func fromDomainSubscriptionChange(change *domain.PaddleSubscriptionChange) *PaddleSubscriptionChangeModel {
	if change == nil {
		return nil
	}
	return &PaddleSubscriptionChangeModel{Action: change.Action, EffectiveAt: change.EffectiveAt}
}

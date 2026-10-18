package mongo

import (
	"fmt"
	paddleD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
)

func linkUpdatedItems(plan *planD.Plan, paddlePlan *paddleD.PaddlePlan) error {
	if len(plan.Items) != len(paddlePlan.Items) {
		return fmt.Errorf("plan and paddle item counts do not match")
	}
	for i, item := range paddlePlan.Items {
		if item == nil || plan.Items[i] == nil || plan.Items[i].ID == "" {
			return fmt.Errorf("plan item at index %d is missing", i)
		}
		item.PlanItemID = plan.Items[i].ID
	}
	return nil
}

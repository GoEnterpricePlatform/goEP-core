package mongo

import (
	"context"
)

func (t *Transaction) deleteLinked(ctx context.Context, id string) error {
	paddlePlan, err := t.PaddlePlanRepo.Find(ctx, id)
	if err != nil {
		return err
	}
	if err := t.PaddlePlanRepo.Delete(ctx, id); err != nil {
		return err
	}
	return t.PlanRepo.Delete(ctx, paddlePlan.PlanID)
}

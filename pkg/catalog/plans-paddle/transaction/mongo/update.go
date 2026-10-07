package mongo

import (
	"context"
	"fmt"

	paddleD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (t *Transaction) Update(ctx context.Context, plan *planD.Plan, paddlePlan *paddleD.PaddlePlan) error {
	session, err := t.Client.StartSession()
	if err != nil {
		return fmt.Errorf("error starting session: %w", err)
	}
	defer session.EndSession(ctx)
	if err := session.StartTransaction(); err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		if err := t.PlanRepo.Update(ctx, plan.ID, plan); err != nil {
			return err
		}
		if err := linkUpdatedItems(plan, paddlePlan); err != nil {
			return err
		}
		if err := t.PaddlePlanRepo.Update(ctx, paddlePlan.ID, paddlePlan); err != nil {
			return err
		}
		return session.CommitTransaction(ctx)
	})
	if err != nil {
		_ = session.AbortTransaction(context.Background())
		return fmt.Errorf("error executing update transaction: %w", err)
	}
	return nil
}

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

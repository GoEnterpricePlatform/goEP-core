package mongo

import (
	"context"
	"fmt"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (t *Transaction) Create(ctx context.Context, plan *domain.Plan, paddlePlan *paddlePlanD.PaddlePlan) error {
	session, err := t.Client.StartSession()
	if err != nil {
		return fmt.Errorf("error starting session: %w", err)
	}

	err = session.StartTransaction()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		err = t.PlanRepo.Insert(ctx, plan)
		if err != nil {
			return err
		}

		paddlePlan.PlanID = plan.ID

		if len(plan.Items) != len(paddlePlan.Items) {
			return fmt.Errorf(
				"plan items count does not match paddle plan items count",
			)
		}

		for i, paddleItem := range paddlePlan.Items {
			if paddleItem == nil {
				continue
			}

			planItem := plan.Items[i]

			if planItem == nil {
				return fmt.Errorf("plan item at index %d is nil", i)
			}

			paddleItem.PlanItemID = planItem.ID
		}

		err = t.PaddlePlanRepo.Insert(ctx, paddlePlan)
		if err != nil {
			return err
		}

		err = session.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if err := session.AbortTransaction(context.Background()); err != nil {
			return fmt.Errorf("error aborting transaction: %w", err)
		}
		return fmt.Errorf("error executing transaction: %w", err)
	}

	return nil

}

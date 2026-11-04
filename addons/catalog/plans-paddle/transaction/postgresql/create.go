package postgresql

import (
	"context"
	"fmt"

	paddleDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"github.com/google/uuid"
)

func (t *Transaction) Create(ctx context.Context, plan *planDomain.Plan, paddlePlan *paddleDomain.PaddlePlan) error {
	if plan == nil || paddlePlan == nil {
		return fmt.Errorf("create Paddle plan: plan and Paddle plan are required")
	}
	if len(plan.Items) != len(paddlePlan.Items) {
		return fmt.Errorf("plan items count does not match Paddle plan items count")
	}
	tx, err := t.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin catalog Paddle create: %w", err)
	}
	defer tx.Rollback(ctx)
	planID := uuid.New()
	if plan.ID != "" {
		planID, err = parseID(plan.ID, "plan")
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO tb_plan(id,name,description,img_path,created_at,updated_at) VALUES($1,$2,$3,$4,COALESCE($5,NOW()),COALESCE($6,NOW()))`, planID, plan.Name, plan.Description, plan.ImgPath, plan.CreatedAt, plan.UpdatedAt)
	if err != nil {
		return catalogError("insert plan", err)
	}
	plan.ID = planID.String()
	paddlePlan.PlanID = plan.ID
	for i, item := range plan.Items {
		if item == nil {
			return fmt.Errorf("plan item at index %d is nil", i)
		}
		paddleItem := paddlePlan.Items[i]
		if paddleItem == nil {
			return fmt.Errorf("Paddle plan item at index %d is nil", i)
		}
		itemID := uuid.New()
		if item.ID != "" {
			itemID, err = parseID(item.ID, "plan item")
			if err != nil {
				return err
			}
		}
		features := item.Features
		if features == nil {
			features = []string{}
		}
		_, err = tx.Exec(ctx, `INSERT INTO tb_plan_item(id,plan_id,img_path,features,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, planID, item.ImgPath, features, item.Status, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return catalogError("insert plan item", err)
		}
		item.ID = itemID.String()
		paddleItem.PlanItemID = item.ID
		for _, rawOptionID := range item.VarOptionIDs {
			optionID, e := parseID(rawOptionID, "variation option")
			if e != nil {
				return e
			}
			if _, e = tx.Exec(ctx, `INSERT INTO tb_plan_item_var_option(plan_item_id,var_option_id) VALUES($1,$2)`, itemID, optionID); e != nil {
				return catalogError("assign plan item option", e)
			}
		}
	}
	paddlePlanID := uuid.New()
	if paddlePlan.ID != "" {
		paddlePlanID, err = parseID(paddlePlan.ID, "Paddle plan")
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO tb_paddle_plan(id,plan_id,paddle_product_id,"order",created_at,updated_at) VALUES($1,$2,$3,$4,COALESCE($5,NOW()),COALESCE($6,NOW()))`, paddlePlanID, planID, paddlePlan.PaddleProductID, paddlePlan.Order, paddlePlan.CreatedAt, paddlePlan.UpdatedAt)
	if err != nil {
		return catalogError("insert Paddle plan", err)
	}
	paddlePlan.ID = paddlePlanID.String()
	for _, item := range paddlePlan.Items {
		itemID := uuid.New()
		if item.ID != "" {
			itemID, err = parseID(item.ID, "Paddle plan item")
			if err != nil {
				return err
			}
		}
		planItemID, e := parseID(item.PlanItemID, "plan item")
		if e != nil {
			return e
		}
		_, e = tx.Exec(ctx, `INSERT INTO tb_paddle_plan_item(id,paddle_plan_id,plan_item_id,paddle_price_id,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, paddlePlanID, planItemID, item.PaddlePriceID, item.Status, item.CreatedAt, item.UpdatedAt)
		if e != nil {
			return catalogError("insert Paddle plan item", e)
		}
		item.ID = itemID.String()
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit catalog Paddle create: %w", err)
	}
	return nil
}

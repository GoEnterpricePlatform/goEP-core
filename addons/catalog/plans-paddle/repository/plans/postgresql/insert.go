package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, paddlePlan *domain.PaddlePlan) error {
	if paddlePlan == nil {
		return fmt.Errorf("insert Paddle plan: nil domain record")
	}
	id := uuid.New()
	if paddlePlan.ID != "" {
		var err error
		id, err = parseID(paddlePlan.ID, "Paddle plan")
		if err != nil {
			return err
		}
	}
	planID, err := parseID(paddlePlan.PlanID, "plan")
	if err != nil {
		return err
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin Paddle plan insert: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO tb_paddle_plan(id,plan_id,paddle_product_id,"order",created_at,updated_at) VALUES($1,$2,$3,$4,COALESCE($5,NOW()),COALESCE($6,NOW()))`, id, planID, paddlePlan.PaddleProductID, paddlePlan.Order, paddlePlan.CreatedAt, paddlePlan.UpdatedAt)
	if err != nil {
		return dbError("insert Paddle plan", err)
	}
	paddlePlan.ID = id.String()
	for _, item := range paddlePlan.Items {
		if item == nil {
			continue
		}
		itemID := uuid.New()
		if item.ID != "" {
			itemID, err = parseID(item.ID, "Paddle plan item")
			if err != nil {
				return err
			}
		}
		planItemID, err := parseID(item.PlanItemID, "plan item")
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO tb_paddle_plan_item(id,paddle_plan_id,plan_item_id,paddle_price_id,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, id, planItemID, item.PaddlePriceID, item.Status, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return dbError("insert Paddle plan item", err)
		}
		item.ID = itemID.String()
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit Paddle plan insert: %w", err)
	}
	return nil
}

package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, rawID string, paddlePlan *domain.PaddlePlan) error {
	id, err := parseID(rawID, "Paddle plan")
	if err != nil {
		return err
	}
	planID, err := parseID(paddlePlan.PlanID, "plan")
	if err != nil {
		return err
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin Paddle plan update: %w", err)
	}
	defer tx.Rollback(ctx)
	res, err := tx.Exec(ctx, `UPDATE tb_paddle_plan SET paddle_product_id=$3,"order"=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1 AND plan_id=$2`, id, planID, paddlePlan.PaddleProductID, paddlePlan.Order, paddlePlan.UpdatedAt)
	if err != nil {
		return dbError("update Paddle plan", err)
	}
	if res.RowsAffected() == 0 {
		return sharedD.ErrNotFound
	}
	if _, err := tx.Exec(ctx, `DELETE FROM tb_paddle_plan_item WHERE paddle_plan_id=$1`, id); err != nil {
		return fmt.Errorf("replace Paddle plan items: %w", err)
	}
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
		linkedPlanItemID, err := parseID(item.PlanItemID, "plan item")
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO tb_paddle_plan_item(id,paddle_plan_id,plan_item_id,paddle_price_id,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, id, linkedPlanItemID, item.PaddlePriceID, item.Status, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return dbError("replace Paddle plan item", err)
		}
		item.ID = itemID.String()
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit Paddle plan update: %w", err)
	}
	return nil
}

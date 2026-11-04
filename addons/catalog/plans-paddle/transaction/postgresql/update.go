package postgresql

import (
	"context"
	"fmt"

	paddleDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	planDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (t *Transaction) Update(ctx context.Context, plan *planDomain.Plan, paddlePlan *paddleDomain.PaddlePlan) error {
	if plan == nil || paddlePlan == nil {
		return fmt.Errorf("update Paddle plan: plan and Paddle plan are required")
	}
	if len(plan.Items) != len(paddlePlan.Items) {
		return fmt.Errorf("plan and Paddle item counts do not match")
	}
	planID, err := parseID(plan.ID, "plan")
	if err != nil {
		return err
	}
	paddlePlanID, err := parseID(paddlePlan.ID, "Paddle plan")
	if err != nil {
		return err
	}
	tx, err := t.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin catalog Paddle update: %w", err)
	}
	defer tx.Rollback(ctx)
	res, err := tx.Exec(ctx, `UPDATE tb_plan SET name=$2,description=$3,img_path=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1`, planID, plan.Name, plan.Description, plan.ImgPath, plan.UpdatedAt)
	if err != nil {
		return catalogError("update plan", err)
	}
	if res.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	keep := make(map[uuid.UUID]struct{}, len(plan.Items))
	for _, item := range plan.Items {
		if item == nil {
			return fmt.Errorf("plan contains a nil item")
		}
		id, e := parseID(item.ID, "plan item")
		if e != nil {
			return e
		}
		keep[id] = struct{}{}
	}
	rows, err := tx.Query(ctx, `SELECT id FROM tb_plan_item WHERE plan_id=$1`, planID)
	if err != nil {
		return fmt.Errorf("list plan items for Paddle update: %w", err)
	}
	var deleteIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if e := rows.Scan(&id); e != nil {
			rows.Close()
			return e
		}
		if _, ok := keep[id]; !ok {
			deleteIDs = append(deleteIDs, id)
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()
	for _, id := range deleteIDs {
		if _, err := tx.Exec(ctx, `DELETE FROM tb_plan_item WHERE id=$1 AND plan_id=$2`, id, planID); err != nil {
			return catalogError("delete plan item", err)
		}
	}
	for _, item := range plan.Items {
		id, _ := uuid.Parse(item.ID)
		features := item.Features
		if features == nil {
			features = []string{}
		}
		res, err := tx.Exec(ctx, `UPDATE tb_plan_item SET img_path=$3,features=$4,updated_at=COALESCE($5,NOW()) WHERE id=$1 AND plan_id=$2`, id, planID, item.ImgPath, features, item.UpdatedAt)
		if err != nil {
			return catalogError("update plan item", err)
		}
		if res.RowsAffected() == 0 {
			return shared.ErrNotFound
		}
		if _, err = tx.Exec(ctx, `DELETE FROM tb_plan_item_var_option WHERE plan_item_id=$1`, id); err != nil {
			return err
		}
		for _, raw := range item.VarOptionIDs {
			optionID, e := parseID(raw, "variation option")
			if e != nil {
				return e
			}
			if _, e = tx.Exec(ctx, `INSERT INTO tb_plan_item_var_option(plan_item_id,var_option_id) VALUES($1,$2)`, id, optionID); e != nil {
				return catalogError("assign plan item option", e)
			}
		}
	}
	res, err = tx.Exec(ctx, `UPDATE tb_paddle_plan SET paddle_product_id=$2,"order"=$3,updated_at=COALESCE($4,NOW()) WHERE id=$1 AND plan_id=$5`, paddlePlanID, paddlePlan.PaddleProductID, paddlePlan.Order, paddlePlan.UpdatedAt, planID)
	if err != nil {
		return catalogError("update Paddle plan", err)
	}
	if res.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	if _, err = tx.Exec(ctx, `DELETE FROM tb_paddle_plan_item WHERE paddle_plan_id=$1`, paddlePlanID); err != nil {
		return err
	}
	for i, item := range paddlePlan.Items {
		if item == nil {
			return fmt.Errorf("Paddle plan item at index %d is nil", i)
		}
		itemID := uuid.New()
		if item.ID != "" {
			itemID, err = parseID(item.ID, "Paddle plan item")
			if err != nil {
				return err
			}
		}
		if i >= len(plan.Items) || plan.Items[i] == nil {
			return fmt.Errorf("linked plan item at index %d is nil", i)
		}
		linkedID, e := parseID(plan.Items[i].ID, "plan item")
		if e != nil {
			return e
		}
		item.PlanItemID = linkedID.String()
		_, e = tx.Exec(ctx, `INSERT INTO tb_paddle_plan_item(id,paddle_plan_id,plan_item_id,paddle_price_id,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),COALESCE($7,NOW()))`, itemID, paddlePlanID, linkedID, item.PaddlePriceID, item.Status, item.CreatedAt, item.UpdatedAt)
		if e != nil {
			return catalogError("replace Paddle plan item", e)
		}
		item.ID = itemID.String()
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit catalog Paddle update: %w", err)
	}
	return nil
}

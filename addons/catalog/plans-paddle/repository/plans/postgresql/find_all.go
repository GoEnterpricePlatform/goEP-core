package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindAll(ctx context.Context, limit, page int64) ([]*domain.PaddlePlan, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	rows, err := r.Pool.Query(ctx, `SELECT id::text,plan_id::text,paddle_product_id,"order",created_at,updated_at FROM tb_paddle_plan ORDER BY "order",created_at DESC,id DESC LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("list Paddle plans: %w", err)
	}
	plans := make([]*domain.PaddlePlan, 0)
	var ids []uuid.UUID
	for rows.Next() {
		p := &domain.PaddlePlan{}
		if err := rows.Scan(&p.ID, &p.PlanID, &p.PaddleProductID, &p.Order, &p.CreatedAt, &p.UpdatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan Paddle plan: %w", err)
		}
		id, _ := uuid.Parse(p.ID)
		ids = append(ids, id)
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate Paddle plans: %w", err)
	}
	rows.Close()
	for i, id := range ids {
		plans[i].Items, err = r.findItems(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return plans, nil
}
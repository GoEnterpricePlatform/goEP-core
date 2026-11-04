package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindAll(ctx context.Context, limit, page int64) ([]*domain.Plan, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	rows, err := r.Pool.Query(ctx, `SELECT id::text,name,description,img_path,created_at,updated_at FROM tb_plan ORDER BY created_at DESC,id LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("find plans: %w", err)
	}
	plans := make([]*domain.Plan, 0)
	for rows.Next() {
		p := &domain.Plan{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImgPath, &p.CreatedAt, &p.UpdatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate plans: %w", err)
	}
	rows.Close()
	for _, p := range plans {
		id, _ := uuid.Parse(p.ID)
		p.Items, err = r.planItems(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return plans, nil
}
package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) Find(ctx context.Context, rawID string) (*domain.PaddlePlan, error) {
	id, err := parseID(rawID, "Paddle plan")
	if err != nil {
		return nil, err
	}
	p := &domain.PaddlePlan{}
	err = r.Pool.QueryRow(ctx, `SELECT id::text,plan_id::text,paddle_product_id,"order",created_at,updated_at FROM tb_paddle_plan WHERE id=$1`, id).Scan(&p.ID, &p.PlanID, &p.PaddleProductID, &p.Order, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sharedD.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find Paddle plan: %w", err)
	}
	p.Items, err = r.findItems(ctx, id)
	if err != nil {
		return nil, err
	}
	return p, nil
}
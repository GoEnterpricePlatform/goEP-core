package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByPaddleProductID(ctx context.Context, paddleProductID string) (*domain.PaddlePlan, error) {
	var id string
	err := r.Pool.QueryRow(ctx, `SELECT id::text FROM tb_catalog_paddle_plan WHERE paddle_product_id=$1`, paddleProductID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sharedD.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find Paddle plan by product ID: %w", err)
	}
	return r.Find(ctx, id)
}

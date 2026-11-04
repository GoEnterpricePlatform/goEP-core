package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.PaddlePlanRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewPlanPaddleRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

func parseID(raw, kind string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s ID: %w", shared.ErrIncorrectID, kind, err)
	}
	return id, nil
}

func dbError(op string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s: %w", shared.ErrDuplicateKey, op, err)
		}
		if pgErr.Code == "23503" {
			return fmt.Errorf("%w: %s: %w", shared.ErrNotFound, op, err)
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

func (r *Repository) findItems(ctx context.Context, paddlePlanID uuid.UUID) ([]*domain.PaddlePlanItem, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id::text,plan_item_id::text,paddle_price_id,status,created_at,updated_at FROM tb_paddle_plan_item WHERE paddle_plan_id=$1 ORDER BY created_at,id`, paddlePlanID)
	if err != nil {
		return nil, fmt.Errorf("find Paddle plan items: %w", err)
	}
	items := make([]*domain.PaddlePlanItem, 0)
	for rows.Next() {
		item := &domain.PaddlePlanItem{}
		if err := rows.Scan(&item.ID, &item.PlanItemID, &item.PaddlePriceID, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan Paddle plan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("iterate Paddle plan items: %w", err)
	}
	rows.Close()
	return items, nil
}

package postgresql

import (
	"context"
	"errors"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) find(ctx context.Context, query string, arg any) (*domain.PaddleCustomer, error) {
	c := &domain.PaddleCustomer{}
	if err := r.Pool.QueryRow(ctx, query, arg).Scan(&c.ID, &c.UserID, &c.PaddleCustomerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedD.ErrNotFound
		}
		return nil, fmt.Errorf("find Paddle customer: %w", err)
	}
	return c, nil
}

package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) Upsert(ctx context.Context, customer *domain.PaddleCustomer) error {
	if customer == nil || customer.UserID == "" || customer.PaddleCustomerID == "" {
		return fmt.Errorf("Paddle customer requires user ID and Paddle customer ID")
	}
	userID, err := uuid.Parse(customer.UserID)
	if err != nil {
		return fmt.Errorf("%w: invalid user ID: %w", shared.ErrIncorrectID, err)
	}
	id := uuid.New()
	if customer.ID != "" {
		id, err = uuid.Parse(customer.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid Paddle customer ID: %w", shared.ErrIncorrectID, err)
		}
	}
	err = r.Pool.QueryRow(ctx, `INSERT INTO tb_paddle_customer(id,user_id,paddle_customer_id) VALUES($1,$2,$3) ON CONFLICT(user_id) DO UPDATE SET paddle_customer_id=EXCLUDED.paddle_customer_id RETURNING id::text`, id, userID, customer.PaddleCustomerID).Scan(&customer.ID)
	if err != nil {
		return fmt.Errorf("upsert Paddle customer: %w", err)
	}
	return nil
}

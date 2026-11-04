package postgresql

import (
	"context"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
)

func (r *Repository) FindByPaddleCustomerID(ctx context.Context, paddleID string) (*domain.PaddleCustomer, error) {
	return r.find(ctx, `SELECT id::text,user_id::text,paddle_customer_id FROM tb_paddle_customer WHERE paddle_customer_id=$1`, paddleID)
}

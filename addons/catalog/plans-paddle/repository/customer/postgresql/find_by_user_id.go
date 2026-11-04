package postgresql

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByUserID(ctx context.Context, rawID string) (*domain.PaddleCustomer, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	return r.find(ctx, `SELECT id::text,user_id::text,paddle_customer_id FROM tb_paddle_customer WHERE user_id=$1`, id)
}

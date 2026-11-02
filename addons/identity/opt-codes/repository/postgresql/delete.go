package postgresql

import (
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/google/uuid"

	"context"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	otpID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: invalid otp ID: %w", sharedD.ErrIncorrectID, err)
	}
	result, err := r.Pool.Exec(ctx, `DELETE FROM tb_identity_otp_code WHERE id=$1`, otpID)
	if err != nil {
		return fmt.Errorf("delete otp code: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: otp %s not found", sharedD.ErrNotFound, id)
	}
	return nil
}

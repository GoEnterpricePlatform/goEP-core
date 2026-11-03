package postgresql

import (
	"context"
	"fmt"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (r *Repository) DeleteByRTokenID(ctx context.Context, tokenID string) error {
	result, err := r.Pool.Exec(ctx, `DELETE FROM tb_identity_session WHERE refresh_token_id=$1`, tokenID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: session not found", sharedD.ErrNotFound)
	}
	return nil
}
package postgresql

import (
	"context"
	"fmt"
)

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var ok bool
	err := r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_identity_user WHERE email=$1)`, email).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check user email existence: %w", err)
	}
	return ok, nil
}
package postgresql

import (
	"context"
	"fmt"
)

func (r *Repository) Exists(ctx context.Context, name string) (bool, error) {
	var ok bool
	err := r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_identity_role WHERE name=$1)`, name).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check role existence: %w", err)
	}
	return ok, nil
}
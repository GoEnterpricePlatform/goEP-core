package postgresql

import "context"

func (r *Repository) HasByVariation(ctx context.Context, rawID string) (bool, error) {
	id, err := parseID(rawID, "variation")
	if err != nil {
		return false, err
	}
	var exists bool
	err = r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tb_var_option WHERE variation_id=$1)`, id).Scan(&exists)
	return exists, err
}
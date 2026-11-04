package postgresql

import "context"

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := r.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM tb_plan`).Scan(&n)
	return n, err
}
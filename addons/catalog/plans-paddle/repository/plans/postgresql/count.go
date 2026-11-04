package postgresql

import "context"

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM tb_paddle_plan`).Scan(&count)
	return count, err
}
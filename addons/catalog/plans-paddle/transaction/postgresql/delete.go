package postgresql

import (
	"context"
	"fmt"

	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
)

func (t *Transaction) Delete(ctx context.Context, rawID string) error {
	id, err := parseID(rawID, "Paddle plan")
	if err != nil {
		return err
	}
	tx, err := t.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin catalog Paddle delete: %w", err)
	}
	defer tx.Rollback(ctx)
	var planID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT plan_id FROM tb_paddle_plan WHERE id=$1`, id).Scan(&planID); err != nil {
		return catalogError("find linked Paddle plan", err)
	}
	res, err := tx.Exec(ctx, `DELETE FROM tb_plan WHERE id=$1`, planID)
	if err != nil {
		return catalogError("delete linked plan", err)
	}
	if res.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit catalog Paddle delete: %w", err)
	}
	return nil
}

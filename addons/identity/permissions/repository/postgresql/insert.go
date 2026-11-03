package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/permissions/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Insert(ctx context.Context, p *domain.Permission) error {
	id := uuid.New()
	if p.ID != "" {
		parsed, err := uuid.Parse(p.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid permission id: %w", sharedD.ErrIncorrectID, err)
		}
		id = parsed
	}
	_, err := r.Pool.Exec(ctx, `INSERT INTO tb_identity_permission (id, name, type) VALUES ($1, $2, $3)`, id, p.Name, p.Type)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: insert permission: %w", sharedD.ErrDuplicateKey, err)
		}
		return fmt.Errorf("insert permission: %w", err)
	}
	p.ID = id.String()
	return nil
}

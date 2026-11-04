package postgresql

import (
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/variations/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.VariationRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewVariationRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

func parseID(id, kind string) (uuid.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s ID: %w", shared.ErrIncorrectID, kind, err)
	}
	return u, nil
}

func dbError(op string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s: %w", shared.ErrDuplicateKey, op, err)
		}
		if pgErr.Code == "23503" {
			return fmt.Errorf("%w: %s: %w", shared.ErrNotFound, op, err)
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

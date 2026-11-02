package postgresql

import (
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.UserRepo = (*Repository)(nil)

type Repository struct { 
	Pool *pgxpool.Pool 
}

func NewUserRepo(pool *pgxpool.Pool) *Repository { 
	return &Repository {
		Pool: pool,
	} 
}

func userError(op string, err error) error {
	var p *pgconn.PgError
	if errors.As(err, &p) && p.Code == "23505" {
		return fmt.Errorf("%w: %s: %w", shared.ErrDuplicateKey, op, err)
	}
	return fmt.Errorf("%s: %w", op, err)
}

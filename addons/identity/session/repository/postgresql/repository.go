package postgresql

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/session/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.SessionRepo = (*Repository)(nil)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Pool: pool,
	}
}

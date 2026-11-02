package postgresql

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.OtpCodeRepo = (*Repository)(nil)

type Repository struct { 
	Pool *pgxpool.Pool 
}

func NewOtpCodeRepo(pool *pgxpool.Pool) *Repository { 
	return &Repository{
		Pool: pool,
	} 
}

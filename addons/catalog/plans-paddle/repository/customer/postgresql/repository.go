package postgresql

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.PaddleCustomerRepo = (*Repository)(nil)

type Repository struct{ Pool *pgxpool.Pool }

func NewPaddleCustomerRepo(pool *pgxpool.Pool) *Repository { return &Repository{Pool: pool} }

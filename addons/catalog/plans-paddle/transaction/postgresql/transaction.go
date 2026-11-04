package postgresql

import (
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/port"
	shared "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ port.PaddlePlanTx = (*Transaction)(nil)

type Transaction struct{ Pool *pgxpool.Pool }

func NewPlanPaddleTx(pool *pgxpool.Pool) *Transaction { return &Transaction{Pool: pool} }

func parseID(raw, kind string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s ID: %w", shared.ErrIncorrectID, kind, err)
	}
	return id, nil
}

func catalogError(op string, err error) error { return fmt.Errorf("%s: %w", op, err) }

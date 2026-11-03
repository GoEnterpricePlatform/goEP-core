package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/domain"
	"github.com/jackc/pgx/v5/pgconn"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"

	"github.com/google/uuid"
)

func (r *Repository) Insert(ctx context.Context, otp *domain.OtpCode) error {
	if otp == nil {
		return fmt.Errorf("insert otp code: nil domain record")
	}
	userID, err := uuid.Parse(otp.UserID)
	if err != nil {
		return fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}

	id := uuid.New()
	if otp.ID != "" {
		id, err = uuid.Parse(otp.ID)
		if err != nil {
			return fmt.Errorf("%w: invalid otp ID: %w", sharedD.ErrIncorrectID, err)
		}
	}
	_, err = r.Pool.Exec(ctx, `INSERT INTO tb_identity_otp_code (id, user_id, otp_code, purpose, used, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, id, userID, otp.OptCode, otp.Purpose, otp.Used, otp.ExpiresAt, otp.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: insert otp code: %w", sharedD.ErrDuplicateKey, err)
		}
		return fmt.Errorf("error inserting otp code: %w", err)
	}
	otp.ID = id.String()
	return nil
}

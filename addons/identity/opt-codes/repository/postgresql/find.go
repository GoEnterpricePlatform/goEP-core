package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/opt-codes/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

func (r *Repository) Find(ctx context.Context, id, userID string) (*domain.OtpCode, error) {
	otpID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid otp ID: %w", sharedD.ErrIncorrectID, err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID: %w", sharedD.ErrIncorrectID, err)
	}
	var otp domain.OtpCode
	var createdAt, expiresAt time.Time
	err = r.Pool.QueryRow(ctx, `SELECT id, user_id, otp_code, purpose, used, expires_at, created_at FROM tb_identity_otp_code WHERE id=$1 AND user_id=$2`, otpID, uid).Scan(&otpID, &uid, &otp.OptCode, &otp.Purpose, &otp.Used, &expiresAt, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: otp %s not found", sharedD.ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("find otp code: %w", err)
	}
	otp.ID, otp.UserID = otpID.String(), uid.String()
	otp.ExpiresAt, otp.CreatedAt = &expiresAt, &createdAt
	return &otp, nil
}

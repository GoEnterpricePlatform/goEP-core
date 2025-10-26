package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/opt-codes/domain"
)

type OtpCodeRepo interface {
	Insert(ctx context.Context, otp *domain.OtpCode) error
	Find(ctx context.Context, id, userID string) (*domain.OtpCode, error)
	Delete(ctx context.Context, id string) error
}

type OtpCodeSrv interface {
	Create(ctx context.Context, otp *domain.OtpCode) error
	Get(ctx context.Context, id, userID string) (*domain.OtpCode, error)
	Delete(ctx context.Context, id string) error
}

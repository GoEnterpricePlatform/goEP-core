package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/features/opt-codes/domain"
)

type OtpCodeRepo interface {
	Insert(ctx context.Context, otp *domain.OtpCode) error
}

type OtpCodeSrv interface {
	Create(ctx context.Context, otp *domain.OtpCode) error
}

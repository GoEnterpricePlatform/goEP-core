package port

import (
	"context"

	"github.com/amorindev/go-cms-tmpl/pkg/features/session/domain"
)

type SessionRepo interface {
	Insert(ctx context.Context, session *domain.Session) error
	DeleteByRTokenID(ctx context.Context, rTokenID string) error
	FindByRTokenID(ctx context.Context, rTokenID, userID string) (*domain.Session, error)
}

type SessionSrv interface {
	Create(ctx context.Context, session *domain.Session, roles []string, email string) error
	DeleteByRTokenID(ctx context.Context, rTokenID string) error
	GetByRTokenID(ctx context.Context, rTokenID, userID string) (*domain.Session, error)
}

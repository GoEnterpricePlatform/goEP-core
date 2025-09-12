package port

import (
	"context"

	"github.com/amorindev/go-tmpl/pkg/app/session/domain"
)

type SessionRepo interface {
	Insert(ctx context.Context, session *domain.Session) error
	DeleteByRTokenID(ctx context.Context, rTokenID string) error
}

type SessionSrv interface {
	Create(ctx context.Context, session *domain.Session, roles []string, email string) error
	DeleteByRTokenID(ctx context.Context, rTokenID string) error
}

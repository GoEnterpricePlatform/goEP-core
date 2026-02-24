package port

import (
	"context"

	sessionD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
)

type AdminSrv interface {
	ExistsAdmin(ctx context.Context) (bool, error)
	SignUp(ctx context.Context, user *domain.User) error
	SignIn(ctx context.Context, email string, password string, rememberMe bool) (*domain.User, *sessionD.Session, error)
}

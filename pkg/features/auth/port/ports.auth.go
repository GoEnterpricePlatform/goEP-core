package port

import (
	"context"

	sessionD "github.com/amorindev/go-tmpl/pkg/features/session/domain"
	userD "github.com/amorindev/go-tmpl/pkg/features/users/domain"
)

type AuthSrv interface {
	SignUp(ctx context.Context, user *userD.User) (string, error)
	SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, error)
}

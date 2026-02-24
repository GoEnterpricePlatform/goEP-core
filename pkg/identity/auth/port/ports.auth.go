package port

import (
	"context"

	sessionD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	userD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
)

type AuthSrv interface {
	SignUp(ctx context.Context, user *userD.User) (string, error)
	SignIn(ctx context.Context, email string, password string, rememberMe bool) (*userD.User, *sessionD.Session, string, error)
	ResendVerifyEmail(ctx context.Context, email string) (string, error)
	VerifyEmail(ctx context.Context, otpID, otpCode, userID string) (*userD.User, *sessionD.Session, error)
	SignOut(ctx context.Context, rTokenID string) error
	RefreshToken(ctx context.Context, rTokenID, userID string) (*sessionD.Session, error)
	GetSession(ctx context.Context, rTokenID, userID string) (*domain.User, *sessionD.Session, error)
}

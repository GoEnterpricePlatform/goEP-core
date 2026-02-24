package port

import "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"

type TokenSrv interface {
	CreateAccessToken(userID string, email string, roles []string, permissions []string) (string, int64, error)
	CreateRefreshToken(userID string, rememberMe bool) (string, string, int64, error)
	ParseAccessToken(tokenStr string) (*claim.AccessTokenClaims, error)
	ParseRefreshToken(tokenStr string) (*claim.RefreshTokenClaims, error)
}

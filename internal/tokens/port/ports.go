package port

import "github.com/amorindev/go-tmpl/internal/tokens/claim"

type TokenSrv interface {
	CreateAccessToken(userID string, email string, roles []string) (string, int64, error)
	CreateRefreshToken(userID string, rememberMe bool) (string, string, int64, error)
	ParseAccessToken(tokenStr string) (*claim.AccessTokenClaims, error)
	ParseRefreshToken(tokenStr string) (*claim.RefreshTokenClaims, error)
}

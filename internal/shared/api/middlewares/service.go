package middlewares

import "github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/port"

type AuthMiddleware struct {
	TokenSrv port.TokenSrv
}

func NewAuthMdw(tokenSrv port.TokenSrv) *AuthMiddleware {
	return &AuthMiddleware{
		TokenSrv: tokenSrv,
	}
}

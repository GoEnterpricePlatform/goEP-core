package middlewares

import "github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/port"

type AuthMiddleware struct {
	TokenSrv port.TokenSrv
}

func NewAuthMdw(tokenSrv port.TokenSrv) *AuthMiddleware {
	return &AuthMiddleware{
		TokenSrv: tokenSrv,
	}
}

package middlewares

import (
	"time"

	authP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/port"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/port"
)

type MdwSrvTmpl struct {
	TokenSrv           tokenP.TokenSrv
	AuthSrv            authP.AuthSrv
	CookieSrv          cookieP.CookieSrv
	JwtAccessCookieDur time.Duration
}

func NewMdwSrvTmpl(
	tokenSrv tokenP.TokenSrv,
	authSrv authP.AuthSrv,
	cookieSrv cookieP.CookieSrv,
) *MdwSrvTmpl {
	return &MdwSrvTmpl{
		TokenSrv:  tokenSrv,
		AuthSrv:   authSrv,
		CookieSrv: cookieSrv,
	}
}

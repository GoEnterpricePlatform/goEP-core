package middlewares

import (
	"time"

	authP "github.com/amorindev/go-cms-tmpl/pkg/identity/auth/port"
	tokenP "github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/port"
	cookieP "github.com/amorindev/go-cms-tmpl/pkg/shared/api/handler/cookie/port"
)

type MdwSrvTmpl struct {
	TokenSrv            tokenP.TokenSrv
	AuthSrv             authP.AuthSrv
	CookieSrv           cookieP.CookieSrv
	JwtAccessCookieDur  time.Duration
}

func NewMdwSrvTmpl(
	tokenSrv tokenP.TokenSrv,
	authSrv authP.AuthSrv,
	cookieSrv cookieP.CookieSrv,
) *MdwSrvTmpl {
	return &MdwSrvTmpl{
		TokenSrv:           tokenSrv,
		AuthSrv:            authSrv,
		CookieSrv:          cookieSrv,
	}
}

package handler

import (
	"net/http"

	adminP "github.com/amorindev/go-cms-tmpl/pkg/identity/admin/port"
	cookieP "github.com/amorindev/go-cms-tmpl/pkg/shared/api/handler/cookie/port"
	"github.com/amorindev/go-cms-tmpl/web/admin/renderer"
	"github.com/amorindev/go-cms-tmpl/web/shared/api/middlewares"
)

type Handler struct {
	AdminSrv      adminP.AdminSrv
	CookieSrv     cookieP.CookieSrv
	ApiBaseUrl    string
	AdminRenderer *renderer.Renderer
	MdwSrvTmpl    *middlewares.MdwSrvTmpl
}

func NewAdminHandler(
	adminSrv adminP.AdminSrv,
	cookieSrv cookieP.CookieSrv,
	apiBaseUrl string,
	adminRenderer *renderer.Renderer,
	mdwSrvtmpl *middlewares.MdwSrvTmpl,
) *Handler {
	h := &Handler{
		ApiBaseUrl:    apiBaseUrl,
		AdminSrv:      adminSrv,
		CookieSrv:     cookieSrv,
		AdminRenderer: adminRenderer,
		MdwSrvTmpl:    mdwSrvtmpl,
	}

	return h
}


func (h Handler) RegisterRoutes(mux *http.ServeMux, muxV1 *http.ServeMux) {
	// Redirect
	mux.HandleFunc("/admin", h.Admin)

	// Pages - render html
	muxV1.Handle("/admin/home", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.HomePage)))
	muxV1.Handle("/admin/other", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.OtherPage)))
	muxV1.HandleFunc("/admin/sign-in", h.SignInPage)
	muxV1.HandleFunc("/admin/sign-up", h.SignUpPage)

	// Actions - form submissions
	muxV1.HandleFunc("POST /admin/auth/sign-up", h.SignUp)
	muxV1.HandleFunc("POST /admin/auth/sign-in", h.SignIn)
}

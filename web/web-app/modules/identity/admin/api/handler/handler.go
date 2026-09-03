package handler

import (
	"net/http"

	adminP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/admin/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/standard-web/admin/renderer"
	cookieP "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
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
	mdwSrvtmpl *middlewares.MdwSrvTmpl,
) *Handler {
	h := &Handler{
		ApiBaseUrl:    apiBaseUrl,
		AdminSrv:      adminSrv,
		CookieSrv:     cookieSrv,
		MdwSrvTmpl:    mdwSrvtmpl,
	}

	return h
}

func (h Handler) RegisterRoutes(mux *http.ServeMux, muxV1 *http.ServeMux) {
	// Redirect
	mux.HandleFunc("/goep-admin", h.GoepAdmin)

	// Pages - render html
	muxV1.HandleFunc("GET /goep-admin/auth/sign-in", h.AdminSignInPage)
	muxV1.HandleFunc("GET /goep-admin/auth/sign-up", h.AdminSignUpPage)
	
	// muxV1.Handle("GET /goep-admin", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.GoepAdminPage)))

	// Actions - form submissions
	// change to goep-admin, and also in the template forms
	muxV1.HandleFunc("POST /admin/auth/sign-up", h.SignUp)
	muxV1.HandleFunc("POST /admin/auth/sign-in", h.SignIn)
}

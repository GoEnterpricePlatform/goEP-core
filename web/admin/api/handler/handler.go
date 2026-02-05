package handler

import (
	"net/http"

	adminP "github.com/amorindev/go-tmpl/pkg/features/admin/port"
	cookieP "github.com/amorindev/go-tmpl/pkg/shared/api/handler/cookie/port"
	"github.com/amorindev/go-tmpl/web/admin/renderer"
)

type Handler struct {
	AdminSrv      adminP.AdminSrv
	CookieSrv     cookieP.CookieSrv
	ApiBaseUrl    string
	AdminRenderer *renderer.Renderer
}

func NewAdminHandler(adminSrv adminP.AdminSrv, cookieSrv cookieP.CookieSrv, apiBaseUrl string, adminRenderer *renderer.Renderer) *Handler {
	h := &Handler{
		ApiBaseUrl:    apiBaseUrl,
		AdminSrv:      adminSrv,
		CookieSrv:     cookieSrv,
		AdminRenderer: adminRenderer,
	}

	return h
}

func (h Handler) RegisterRoutes(mux *http.ServeMux, muxV1 *http.ServeMux) {
	// Redirect
	mux.HandleFunc("/admin", h.Admin)

	// Pages - render html
	muxV1.HandleFunc("/admin/home", h.HomePage)
	muxV1.HandleFunc("/admin/other", h.OtherPage)
	muxV1.HandleFunc("/admin/sign-in", h.SignInPage)
	muxV1.HandleFunc("/admin/sign-up", h.SignUpPage)

	// Actions - form submissions
	muxV1.HandleFunc("POST /admin/auth/sign-up", h.SignUp)
	muxV1.HandleFunc("POST /admin/auth/sign-in", h.SignIn)
}

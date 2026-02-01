package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/admin/port"
	"github.com/amorindev/go-tmpl/web/admin/renderer"
)

type Handler struct {
	AdminSrv      port.AdminSrv
	ApiBaseUrl    string
	AdminRenderer *renderer.Renderer
}

func NewAdminHandler(adminSrv port.AdminSrv, apiBaseUrl string, adminRenderer *renderer.Renderer) *Handler {
	h := &Handler{
		ApiBaseUrl:    apiBaseUrl,
		AdminSrv:      adminSrv,
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
}

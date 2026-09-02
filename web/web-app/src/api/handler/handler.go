package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

type Handler struct {
	PostSrv    port.PostSrv
	MdwSrvTmpl *middlewares.MdwSrvTmpl
}

func NewSrcHandler(mux *http.ServeMux, templateV1 *http.ServeMux, postSrv port.PostSrv, mdwSrvtmpl *middlewares.MdwSrvTmpl) *Handler {
	h := &Handler{
		PostSrv:    postSrv,
		MdwSrvTmpl: mdwSrvtmpl,
	}

	// Redirects requests from "/" to the landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/landing", http.StatusFound)
	})

	mux.HandleFunc("/landing", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)
	mux.HandleFunc("/blog", h.BlogPage)
	mux.HandleFunc("/contact", h.ContactPage)

	templateV1.Handle("/goep-admin/general", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.GeneralPage)))
	templateV1.Handle("/goep-admin/users", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.UsersPage)))
	templateV1.Handle("/goep-admin/roles-permissions", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.RolesPermissionsPage)))
	templateV1.Handle("/goep-admin/settings", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.SettingsPage)))

	return h
}

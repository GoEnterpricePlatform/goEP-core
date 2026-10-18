package handler

import (
	"net/http"

	roleP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
)

type Handler struct {
	PostSrv    port.PostSrv
	RoleSrv    roleP.RoleSrv
	MdwSrvTmpl *middlewares.MdwSrvTmpl
}

func NewSrcHandler(mux *http.ServeMux, templateV1 *http.ServeMux, postSrv port.PostSrv, mdwSrvtmpl *middlewares.MdwSrvTmpl, roleSrv roleP.RoleSrv) *Handler {
	h := &Handler{
		PostSrv:    postSrv,
		RoleSrv:    roleSrv,
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

	templateV1.Handle("/goep-admin/general",
		h.MdwSrvTmpl.Authenticate(
			h.MdwSrvTmpl.RequirePermission(
				"view.general",
				http.HandlerFunc(h.GeneralPage),
			),
		),
	)
	templateV1.Handle("/goep-admin/users",
		h.MdwSrvTmpl.Authenticate(
			h.MdwSrvTmpl.RequirePermission(
				"view.users",
				http.HandlerFunc(h.UsersPage),
			),
		),
	)
	templateV1.Handle("/goep-admin/roles-permissions",
		h.MdwSrvTmpl.Authenticate(
			h.MdwSrvTmpl.RequirePermission(
				"view.roles.permissions",
				http.HandlerFunc(h.RolesPermissionsPage),
			),
		),
	)

	templateV1.Handle("/goep-admin/settings",
		h.MdwSrvTmpl.Authenticate(
			h.MdwSrvTmpl.RequirePermission(
				"view.settings",
				http.HandlerFunc(h.SettingsPage),
			),
		),
	)

	templateV1.HandleFunc("/goep-admin/access-denied", h.AccessDeniedPage)

	return h
}

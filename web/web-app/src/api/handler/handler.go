package handler

import (
	"net/http"
	"path"
	"strings"

	roleP "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	"github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	moduleRegistry "github.com/GoEnterpricePlatform/goEP-core/internal/modules"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/resources"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/src/ui/pages"
)

type Handler struct {
	PostSrv          port.PostSrv
	RoleSrv          roleP.RoleSrv
	MdwSrvTmpl       *middlewares.MdwSrvTmpl
	AvailableModules []moduleRegistry.Manifest
}

func NewSrcHandler(mux *http.ServeMux, templateV1 *http.ServeMux, postSrv port.PostSrv, mdwSrvtmpl *middlewares.MdwSrvTmpl, roleSrv roleP.RoleSrv, availableModules []moduleRegistry.Manifest) *Handler {
	h := &Handler{
		PostSrv:          postSrv,
		RoleSrv:          roleSrv,
		MdwSrvTmpl:       mdwSrvtmpl,
		AvailableModules: availableModules,
	}

	// Redirects requests from "/" to the landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/landing", http.StatusFound)
	})

	mux.HandleFunc("/landing", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)
	mux.HandleFunc("/blog", h.BlogPage)
	mux.HandleFunc("/contact", h.ContactPage)
	templateV1.Handle("GET /goep-admin/modules", h.MdwSrvTmpl.Authenticate(http.HandlerFunc(h.ModulesPage)))

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

func (h *Handler) ModulesPage(w http.ResponseWriter, r *http.Request) {
	cards := make([]pages.ModuleCard, 0, len(h.AvailableModules))
	for _, manifest := range h.AvailableModules {
		iconURL := ""
		if manifest.Icon != "" {
			iconPath := path.Join("modules", manifest.Name, strings.TrimPrefix(manifest.Icon, "static/"))
			iconURL = resources.StaticPath(iconPath)
		}
		cards = append(cards, pages.ModuleCard{
			Name:    manifest.Name,
			Version: manifest.Version,
			Summary: manifest.Summary,
			IconURL: iconURL,
		})
	}

	if err := pages.ModulesPage(cards).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

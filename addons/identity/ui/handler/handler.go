package handler

import (
	"net/http"

	rolePort "github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/port"
	rolesHandler "github.com/GoEnterpricePlatform/goEP-core/addons/identity/ui/roles/handler"
	sharedNavigation "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
	"github.com/a-h/templ"
)

// DashboardLayout lets the web application wrap module content in its shared dashboard shell.
type DashboardLayout func(activePage string, roles, permissions []string, menu []sharedNavigation.Item, content templ.Component) templ.Component

type Handler struct {
	RoleSrv         rolePort.RoleSrv
	MdwSrvTmpl      *middlewares.MdwSrvTmpl
	DashboardLayout DashboardLayout
	Roles           *rolesHandler.Handler
}

func NewHandler(
	roleSrv rolePort.RoleSrv,
	apiBaseURL string,
	mdwSrvTmpl *middlewares.MdwSrvTmpl,
	dashboardLayout DashboardLayout,
) *Handler {
	return &Handler{
		RoleSrv:         roleSrv,
		MdwSrvTmpl:      mdwSrvTmpl,
		DashboardLayout: dashboardLayout,
		Roles:           rolesHandler.NewRolesTmplHandler(roleSrv, apiBaseURL, mdwSrvTmpl),
	}
}

// RegisterRoutes attaches Identity's dashboard pages and role actions to the web template mux.
func (h *Handler) RegisterRoutes(templateV1 *http.ServeMux) {
	templateV1.Handle("/goep-admin/general", h.MdwSrvTmpl.Authenticate(
		h.MdwSrvTmpl.RequirePermission("view.general", http.HandlerFunc(h.GeneralPage)),
	))
	templateV1.Handle("/goep-admin/users", h.MdwSrvTmpl.Authenticate(
		h.MdwSrvTmpl.RequirePermission("view.users", http.HandlerFunc(h.UsersPage)),
	))
	templateV1.Handle("/goep-admin/roles-permissions", h.MdwSrvTmpl.Authenticate(
		h.MdwSrvTmpl.RequirePermission("view.roles.permissions", http.HandlerFunc(h.RolesPermissionsPage)),
	))
	templateV1.Handle("/goep-admin/settings", h.MdwSrvTmpl.Authenticate(
		h.MdwSrvTmpl.RequirePermission("view.settings", http.HandlerFunc(h.SettingsPage)),
	))

	h.Roles.RegisterRoutes(nil, templateV1)
}

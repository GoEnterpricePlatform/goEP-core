package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/claim"
	identityNavigation "github.com/GoEnterpricePlatform/goEP-core/addons/identity/ui/navigation"
	"github.com/a-h/templ"
)

func (h *Handler) renderDashboard(w http.ResponseWriter, r *http.Request, claims *claim.AccessTokenClaims, activePage string, content templ.Component) {
	if claims == nil {
		http.Redirect(w, r, "/v1/goep-admin/auth/sign-in", http.StatusFound)
		return
	}

	if err := h.DashboardLayout(activePage, claims.Roles, claims.Permissions, identityNavigation.AdminMenu, content).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/claim"
	identityPages "github.com/GoEnterpricePlatform/goEP-core/addons/identity/ui/pages"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
)

func (h *Handler) GeneralPage(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(
		middlewares.AccessTokenClaimsTmplIDKey,
	).(*claim.AccessTokenClaims)

	if !ok || claims == nil {
		http.Redirect(
			w,
			r,
			"/v1/goep-admin/auth/sign-in",
			http.StatusFound,
		)
		return
	}

	h.renderDashboard(w, r, claims, "general", identityPages.GeneralContent())
}

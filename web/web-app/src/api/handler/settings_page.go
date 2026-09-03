package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/ui/pages"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
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

	err := pages.SettingsPage(claims).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

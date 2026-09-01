package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/modules/identity/ui/pages"
)

func (h *Handler) AdminSignInPage(w http.ResponseWriter, r *http.Request) {
	err := pages.AdminSignInPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
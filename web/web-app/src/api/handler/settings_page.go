package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/ui/pages"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	err := pages.SettingsPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

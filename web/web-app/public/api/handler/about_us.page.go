package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/ui/pages"
)

func (h *Handler) AboutUsPage(w http.ResponseWriter, r *http.Request) {
	err := pages.AboutUsPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
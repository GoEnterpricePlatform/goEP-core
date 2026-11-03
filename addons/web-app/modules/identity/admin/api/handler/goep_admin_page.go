package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/web-app/src/ui/pages"
)

func (h *Handler) GoepAdminPage(w http.ResponseWriter, r *http.Request) {
	err := pages.GoepAdminPage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

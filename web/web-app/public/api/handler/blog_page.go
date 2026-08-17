package handler

import (
	"net/http"

	mocksdata "github.com/GoEnterpricePlatform/goEP-core/web/web-app/mocks_data"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/ui/pages"
)

func (h *Handler) BlogPage(w http.ResponseWriter, r *http.Request) {
	err := pages.BlogPage(mocksdata.Posts).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/web/shared/templates"
)

// SignInPage renders the admin sign-in page.
func (h Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	h.AdminRenderer.Render(w, "sign-in", templates.ErrorData{})
}

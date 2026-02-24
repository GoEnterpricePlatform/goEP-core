package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/shared/templates"
)

func (h Handler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	h.AdminRenderer.Render(w, "sign-up", templates.ErrorData{})
}

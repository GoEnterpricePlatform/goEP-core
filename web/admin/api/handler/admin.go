package handler

import (
	"net/http"

	"github.com/amorindev/go-cms-tmpl/web/shared/core"
	"github.com/amorindev/go-cms-tmpl/web/shared/templates"
)

// Admin checks if an admin already exists and redirects
// to sign-in or sign-up accordingly. Renders an error page on failure.
func (h Handler) Admin(w http.ResponseWriter, r *http.Request) {
	exists, err := h.AdminSrv.ExistsAdmin(r.Context())
	if err != nil {
		h.AdminRenderer.Render(w, "error", templates.ErrorData{
			ErrorMsg: core.UiErrorResp(err),
		})
	}
	if exists {
		http.Redirect(w, r, "/v1/admin/sign-in", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/v1/admin/sign-up", http.StatusFound)
}

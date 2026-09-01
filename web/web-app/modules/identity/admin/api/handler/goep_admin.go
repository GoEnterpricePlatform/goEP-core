package handler

import (
	"net/http"
)

// GoepAdmin checks if an admin user already exists and redirects
// to admin sign-in or admin sign-up accordingly. Renders an error page on failure.
func (h Handler) GoepAdmin(w http.ResponseWriter, r *http.Request) {
	exists, err := h.AdminSrv.ExistsAdmin(r.Context())
	if err != nil {
		http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
		return
	}
	if exists {
		http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/v1/admin/auth/sign-up", http.StatusFound)
}

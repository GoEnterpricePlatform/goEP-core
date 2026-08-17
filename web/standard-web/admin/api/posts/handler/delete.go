package handler

import (
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Redirect(w, r, "/v1/admin/posts", http.StatusSeeOther)
		return
	}

	err := h.PostSrv.Delete(r.Context(), id)
	if err != nil {
		posts, err := h.PostSrv.GetAll(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.AdminRenderer.Render(w, "base", RespData{
			ActivePage:  "post",
			Posts:       posts,
			DeleteErr:   sharedC.UiErrorResp(err),
			DeleteErrID: id,
		})
		return
	}

	http.Redirect(w, r, "/v1/admin/post", http.StatusSeeOther)
}

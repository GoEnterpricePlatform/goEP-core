package handler

import (
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	"github.com/amorindev/go-cms-tmpl/web/shared/core"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	desc := r.FormValue("desc")

	post := domain.Post{
		ID:    id,
		Title: r.FormValue("title"),
		Desc:  &desc,
	}

	err := h.PostSrv.Update(r.Context(), id, &post)
	if err != nil {
		posts, err := h.PostSrv.GetAll(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.AdminRenderer.Render(w, "base", RespData{
			ActivePage:  "post",
			Posts:       posts,
			EditingID:   id,
			UpdateErr:   core.UiErrorResp(err),
			UpdateErrID: id,
		})
		return
	}

	http.Redirect(w, r, "/v1/admin/post", http.StatusSeeOther)
}

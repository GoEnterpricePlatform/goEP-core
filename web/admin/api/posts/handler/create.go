package handler

import (
	"context"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedC "github.com/amorindev/go-cms-tmpl/web/shared/core"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	desc := r.FormValue("desc")
	p := domain.Post{
		Title: r.FormValue("title"),
		Desc:  &desc,
	}

	err := h.PostSrv.Create(context.Background(), &p)
	if err != nil {
		posts, err := h.PostSrv.GetAll(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		h.AdminRenderer.Render(w, "base", RespData{
			ActivePage: "post",
			Posts:      posts,
			CreateErr:  sharedC.UiErrorResp(err),
		})
		return
	}
	http.Redirect(w, r, "/v1/admin/post", http.StatusSeeOther)
}

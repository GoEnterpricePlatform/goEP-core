package handler

import (
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
)

func (h Handler) PostPage(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostSrv.GetAll(r.Context())

	data := RespData{
		ActivePage: "post",
		ErrorMsg:   "",
		Posts:      posts,
		EditingID:  "",
	}

	if err != nil {
		data.PageErr = sharedC.UiErrorResp(err)
	}
	h.AdminRenderer.Render(w, "base", data)
}

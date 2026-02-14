package handler

import (
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/port"
)

type Handler struct {
	PostSrv port.PostSrv
}

func NewPostApiHandler(muxV1 *http.ServeMux, postSrv port.PostSrv) *Handler {
	h := &Handler{
		PostSrv: postSrv,
	}

	muxV1.HandleFunc("POST /posts", h.Create)
	muxV1.HandleFunc("GET /posts", h.GetAll)
	muxV1.HandleFunc("PUT /posts/{id}", h.Update)
	muxV1.HandleFunc("PATCH /posts/{id}", h.Patch)
	muxV1.HandleFunc("DELETE /posts/{id}", h.Delete)

	return h
}

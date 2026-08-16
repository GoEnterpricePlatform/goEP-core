package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	PostSrv    port.PostSrv
	AuthApiMdw *middlewares.AuthMiddleware
}

func NewPostApiHandler(muxV1 *http.ServeMux, postSrv port.PostSrv, authApiMdw *middlewares.AuthMiddleware) *Handler {
	h := &Handler{
		PostSrv: postSrv,
	}

	muxV1.Handle("POST /posts", h.AuthApiMdw.AccessTokenMdw(h.Create))
	muxV1.HandleFunc("GET /posts/{id}", h.Get)
	muxV1.HandleFunc("GET /posts", h.GetAll)
	muxV1.Handle("PUT /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Update))
	muxV1.Handle("PATCH /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Patch))
	muxV1.Handle("DELETE /posts/{id}", h.AuthApiMdw.AccessTokenMdw(h.Delete))
	muxV1.HandleFunc("GET /posts/search", h.Search)

	return h
}

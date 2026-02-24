package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/web/admin/renderer"
	postP "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
)

type Handler struct {
	PostSrv       postP.PostSrv
	ApiBaseUrl    string
	AdminRenderer *renderer.Renderer
}

func NewPostWebHandler(postSrv postP.PostSrv, apiBaseUrl string, adminRenderer *renderer.Renderer) *Handler {
	h := &Handler{
		ApiBaseUrl:    apiBaseUrl,
		PostSrv:       postSrv,
		AdminRenderer: adminRenderer,
	}

	return h
}

func (h Handler) RegisterRoutes(muxV1 *http.ServeMux) {

	// Pages - render html
	muxV1.HandleFunc("/admin/post", h.PostPage)

	// Actions - form submissions
	muxV1.HandleFunc("POST /admin/posts", h.Create)
	muxV1.HandleFunc("POST /admin/posts/update", h.Update)
	muxV1.HandleFunc("POST /admin/posts/delete", h.Delete)
	muxV1.HandleFunc("GET /admin/posts/edit", h.Edit)
}

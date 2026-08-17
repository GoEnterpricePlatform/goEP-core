package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
)

type Handler struct {
	ApiBaseUrl string
	PostSrv    port.PostSrv
}

func NewPublicHandler(mux *http.ServeMux, apiBaseUrl string, postSrv port.PostSrv) *Handler {
	h := &Handler{
		ApiBaseUrl: apiBaseUrl,
		PostSrv:    postSrv,
	}

	mux.HandleFunc("/", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)

	return h
}

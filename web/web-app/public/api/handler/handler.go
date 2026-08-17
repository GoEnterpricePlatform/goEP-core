package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
)

type Handler struct {
	PostSrv port.PostSrv
}

func NewPublicHandler(mux *http.ServeMux, postSrv port.PostSrv) *Handler {
	h := &Handler{
		PostSrv: postSrv,
	}

	// Redirects requests from "/" to the landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/landing", http.StatusFound)
	})

	mux.HandleFunc("/landing", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)
	mux.HandleFunc("/blog", h.BlogPage)
	mux.HandleFunc("/contact", h.ContactPage)

	return h
}

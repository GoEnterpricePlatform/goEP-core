package handler

import "net/http"

type Handler struct {
	ApiBaseUrl string
}

func NewPublicHandler(mux *http.ServeMux, apiBaseUrl string) *Handler {
	h := &Handler{
		ApiBaseUrl: apiBaseUrl,
	}

	mux.HandleFunc("/", h.LandingPage)
	mux.HandleFunc("/about_us", h.AboutUsPage)

	return h
}
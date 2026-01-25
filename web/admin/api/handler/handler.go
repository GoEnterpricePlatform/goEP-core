package handler

import "net/http"

type Handler struct {
	ApiBaseUrl string
}

func NewAdminHandler(mux *http.ServeMux, apiBaseUrl string) *Handler {
	h := &Handler{
		ApiBaseUrl: apiBaseUrl,
	}

	mux.HandleFunc("/admin/home", h.HomePage)
	mux.HandleFunc("/admin/other", h.OtherPage)

	return h
}

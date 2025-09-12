package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/app/auth-methods/port"
)

type Handler struct {
	AuthMethodSrv port.AuthMethodSrv
}

func NewAuthMethodHandler(server *http.ServeMux, authMethodSrv port.AuthMethodSrv) *Handler {
	h := &Handler{
		AuthMethodSrv: authMethodSrv,
	}

	server.HandleFunc("POST /auth/sign-up", h.SignUp)
	server.HandleFunc("POST /auth/sign-in", h.SignIn)

	return h
}

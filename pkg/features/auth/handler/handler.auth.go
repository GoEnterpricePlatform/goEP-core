package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/auth/port"
)

type Handler struct {
	AuthSrv port.AuthSrv
}

func NewAuthHandler(server *http.ServeMux, authSrv port.AuthSrv) *Handler {
	h := &Handler{
		AuthSrv: authSrv,
	}

	server.HandleFunc("POST /auth/sign-up", h.SignUp)
	server.HandleFunc("POST /auth/sign-in", h.SignIn)
	server.HandleFunc("POST /auth/resend-verify-email", h.ResendVerifyEmail)

	return h
}
